// Package nats 管理 NATS 连接：未配置 url 时嵌入 nats-server，否则外连。
package nats

import (
	"fmt"
	"log/slog"
	"time"

	natsserver "github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

// Conn 持有客户端连接及可选的嵌入服务。
type Conn struct {
	NC     *nats.Conn
	server *natsserver.Server
}

// Connect 按 url 连接；url 为空则启动嵌入 NATS。
func Connect(url string) (*Conn, error) {
	if url == "" {
		return startEmbedded()
	}
	nc, err := nats.Connect(url,
		nats.Name("lite-ai-gateway"),
		nats.Timeout(5*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("connect nats %s: %w", url, err)
	}
	slog.Info("connected to nats", "url", url)
	return &Conn{NC: nc}, nil
}

func startEmbedded() (*Conn, error) {
	opts := &natsserver.Options{
		Host:      "127.0.0.1",
		Port:      -1,
		JetStream: true,
		StoreDir:  "data/nats",
		NoLog:     true,
		NoSigs:    true,
	}
	srv, err := natsserver.NewServer(opts)
	if err != nil {
		return nil, fmt.Errorf("create embedded nats: %w", err)
	}
	go srv.Start()
	if !srv.ReadyForConnections(10 * time.Second) {
		srv.Shutdown()
		return nil, fmt.Errorf("embedded nats not ready")
	}
	nc, err := nats.Connect(srv.ClientURL(), nats.Name("lite-ai-gateway"))
	if err != nil {
		srv.Shutdown()
		return nil, fmt.Errorf("connect embedded nats: %w", err)
	}
	slog.Info("embedded nats started", "url", srv.ClientURL())
	return &Conn{NC: nc, server: srv}, nil
}

// Close 关闭连接与嵌入服务。
func (c *Conn) Close() {
	if c == nil {
		return
	}
	if c.NC != nil {
		c.NC.Close()
	}
	if c.server != nil {
		c.server.Shutdown()
	}
}
