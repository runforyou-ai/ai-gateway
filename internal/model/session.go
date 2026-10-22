package model

import (
	"time"

	"github.com/uptrace/bun"
)

// Session 对应登录会话。
type Session struct {
	bun.BaseModel `bun:"table:sessions,alias:s"`

	ID        int64     `bun:"id,pk,autoincrement"`
	UserID    int64     `bun:"user_id,notnull"`
	TokenHash string    `bun:"token_hash,notnull"`
	ExpiresAt time.Time `bun:"expires_at,notnull"`
	CreatedAt time.Time `bun:"created_at,nullzero,default:now()"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,default:now()"`
}
