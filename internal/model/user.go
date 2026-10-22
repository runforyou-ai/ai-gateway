package model

import (
	"time"

	"github.com/runforyou-ai/lite-ai-gateway/internal/domain"
	"github.com/uptrace/bun"
)

// User 对应用户账号。
type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID           int64        `bun:"id,pk,autoincrement"`
	Email        string       `bun:"email,notnull"`
	PasswordHash string       `bun:"password_hash,notnull"`
	Role         domain.Role  `bun:"role,notnull"`
	CreatedAt    time.Time    `bun:"created_at,nullzero,default:now()"`
	UpdatedAt    time.Time    `bun:"updated_at,nullzero,default:now()"`
}
