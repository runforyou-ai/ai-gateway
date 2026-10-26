// Package app 组装进程生命周期：数据库、HTTP、信号。
package app

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/runforyou-ai/ai-gateway/internal/api"
	"github.com/runforyou-ai/ai-gateway/internal/auth"
	"github.com/runforyou-ai/ai-gateway/internal/config"
	"github.com/runforyou-ai/ai-gateway/internal/direct"
	"github.com/runforyou-ai/ai-gateway/internal/email"
	"github.com/runforyou-ai/ai-gateway/internal/logging"
	"github.com/runforyou-ai/ai-gateway/internal/migrate"
	"github.com/runforyou-ai/ai-gateway/internal/process"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type sessionAuthorizer struct {
	inner *auth.SessionAuthorizer
}

func (a sessionAuthorizer) Authenticate(ctx context.Context, r direct.CookieRequest) (auth.Principal, error) {
	return a.inner.Authenticate(ctx, r.CookieHeader)
}

// Run 加载配置、迁移、启动服务并阻塞至停止信号。
func Run(cfgPath, pidPath string) error {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return err
	}
	if pid, err := process.RunningPID(pidPath); err != nil {
		return err
	} else if pid > 0 {
		return fmt.Errorf("already running (pid %d); stop it first", pid)
	}

	closeLog, err := logging.Setup(logging.Options{Dir: cfg.LogDir()})
	if err != nil {
		return err
	}
	defer func() { _ = closeLog() }()

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.PostgresDSN())))
	sqldb.SetMaxOpenConns(10)
	sqldb.SetMaxIdleConns(5)
	sqldb.SetConnMaxLifetime(time.Hour)
	if err := sqldb.Ping(); err != nil {
		return fmt.Errorf("ping postgres: %w", err)
	}
	defer sqldb.Close()

	ctx := context.Background()
	if err := migrate.Up(ctx, sqldb); err != nil {
		return err
	}
	db := bun.NewDB(sqldb, pgdialect.New())

	sessionAuth := &auth.SessionAuthorizer{DB: db}
	direct.SetAuthorizer(sessionAuthorizer{inner: sessionAuth})
	mailer := email.New(email.Options{
		From:         cfg.Email.From,
		FromName:     cfg.Email.FromName,
		ResendAPIKey: cfg.Email.Resend.APIKey,
	})
	if cfg.Email.Configured() {
		slog.Info("email enabled", "provider", "resend")
	} else {
		slog.Info("email disabled: from or resend api_key not set")
	}
	backend := direct.NewBackend(db, mailer, cfg.Server.PublicURL)
	handler := api.NewMux(backend, api.MuxOptions{
		DB:           db,
		CookieSecure: strings.HasPrefix(cfg.Server.PublicURL, "https://"),
		Mailer:       mailer,
		PublicURL:    cfg.Server.PublicURL,
	})
	server := &http.Server{
		Addr:              cfg.Addr(),
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}

	ln, err := net.Listen("tcp", cfg.Addr())
	if err != nil {
		return fmt.Errorf("listen %s: %w", cfg.Addr(), err)
	}
	if err := process.WritePID(pidPath); err != nil {
		_ = ln.Close()
		return fmt.Errorf("write pid: %w", err)
	}
	defer func() { _ = process.RemovePID(pidPath) }()

	errCh := make(chan error, 1)
	go func() {
		slog.Info("http server listening", "addr", cfg.Addr(), "env", cfg.Env)
		if err := server.Serve(ln); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer signal.Stop(sigCh)

	for {
		select {
		case err := <-errCh:
			return err
		case sig := <-sigCh:
			switch sig {
			case syscall.SIGHUP:
				slog.Info("reload signal received")
				newCfg, err := config.Load(cfgPath)
				if err != nil {
					slog.Error("reload config failed", "err", err)
					continue
				}
				cfg = newCfg
				slog.Info("config reloaded", "env", cfg.Env, "log_dir", cfg.LogDir())
			default:
				slog.Info("shutdown signal received", "signal", sig.String())
				shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
				defer cancel()
				return server.Shutdown(shutdownCtx)
			}
		}
	}
}
