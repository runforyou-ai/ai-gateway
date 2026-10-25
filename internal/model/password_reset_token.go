package model

import (
	"time"

	"github.com/uptrace/bun"
)

// PasswordResetToken 对应密码重置令牌。
type PasswordResetToken struct {
	bun.BaseModel `bun:"table:password_reset_tokens,alias:prt"`

	ID        int64      `bun:"id,pk,autoincrement"`
	CreatedAt time.Time  `bun:"created_at,nullzero,default:now()"`
	UpdatedAt time.Time  `bun:"updated_at,nullzero,default:now()"`
	UserID    int64      `bun:"user_id,notnull"`
	TokenHash string     `bun:"token_hash,notnull"`
	ExpiresAt time.Time  `bun:"expires_at,notnull"`
	UsedAt    *time.Time `bun:"used_at"`
}
