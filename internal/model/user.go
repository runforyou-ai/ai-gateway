package model

import (
	"time"

	"github.com/runforyou-ai/lite-ai-gateway/internal/domain"
	"github.com/uptrace/bun"
)

// User 对应用户账号。
type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID        int64       `bun:"id,pk,autoincrement"`
	Name      string      `bun:"name,notnull"`
	Email     string      `bun:"email,notnull"`
	Password  string      `bun:"password,notnull"`
	Role      domain.Role `bun:"role,notnull"`
	CreatedAt time.Time   `bun:"created_at,nullzero,default:now()"`
	UpdatedAt time.Time   `bun:"updated_at,nullzero,default:now()"`
}
