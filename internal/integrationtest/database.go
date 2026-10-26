// Package integrationtest 对 actions 做真实 Postgres 集成测试。
//
// 需设置 TEST_POSTGRES_*（至少 TEST_POSTGRES_HOST）；库名须以 _test 结尾。
// 未设置时测试 Skip，不占用普通 go test。
package integrationtest

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/runforyou-ai/ai-gateway/internal/migrate"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var (
	dbOnce sync.Once
	shared *bun.DB
	dbErr  error
)

// Open 返回已迁移的测试库连接；每个用例开始时会 Reset 业务表。
// TEST_POSTGRES_HOST 未设置时 Skip。
func Open(t *testing.T) *bun.DB {
	t.Helper()
	dbOnce.Do(func() {
		shared, dbErr = openOnce()
	})
	if dbErr != nil {
		if skip, ok := dbErr.(errSkip); ok {
			t.Skip(skip.Error())
		}
		t.Fatal(dbErr)
	}
	Reset(t, shared)
	return shared
}

func openOnce() (*bun.DB, error) {
	host := strings.TrimSpace(os.Getenv("TEST_POSTGRES_HOST"))
	if host == "" {
		return nil, errSkip{reason: "TEST_POSTGRES_HOST not set"}
	}
	port := 5432
	if v := strings.TrimSpace(os.Getenv("TEST_POSTGRES_PORT")); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("TEST_POSTGRES_PORT: %w", err)
		}
		port = n
	}
	user := envOr("TEST_POSTGRES_USER", "ai-gateway")
	password := envOr("TEST_POSTGRES_PASSWORD", "ai-gateway")
	database := envOr("TEST_POSTGRES_DATABASE", "ai-gateway_test")
	sslmode := envOr("TEST_POSTGRES_SSLMODE", "disable")
	if !strings.HasSuffix(database, "_test") {
		return nil, fmt.Errorf("TEST_POSTGRES_DATABASE %q must end with _test", database)
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		user, password, host, port, database, sslmode,
	)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	sqldb.SetMaxOpenConns(5)
	sqldb.SetMaxIdleConns(2)
	sqldb.SetConnMaxLifetime(time.Hour)
	if err := sqldb.Ping(); err != nil {
		_ = sqldb.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	ctx := context.Background()
	if err := migrate.Up(ctx, sqldb); err != nil {
		_ = sqldb.Close()
		return nil, err
	}
	return bun.NewDB(sqldb, pgdialect.New()), nil
}

// Reset 清空业务表，保留 goose 迁移历史。
func Reset(t *testing.T, db *bun.DB) {
	t.Helper()
	ctx := context.Background()
	_, err := db.ExecContext(ctx, `
		TRUNCATE TABLE virtual_keys, user_invites, password_reset_tokens, sessions, users, settings, state RESTART IDENTITY CASCADE`)
	if err != nil {
		t.Fatalf("reset tables: %v", err)
	}
}

func envOr(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

type errSkip struct {
	reason string
}

func (e errSkip) Error() string { return e.reason }
