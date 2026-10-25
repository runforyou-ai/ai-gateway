package model

import (
	"time"

	"github.com/runforyou-ai/ai-gateway/internal/domain"
	"github.com/uptrace/bun"
)

// UserInvite 对应用户邀请令牌。
type UserInvite struct {
	bun.BaseModel `bun:"table:user_invites,alias:ui"`

	ID        int64       `bun:"id,pk,autoincrement"`
	CreatedAt time.Time   `bun:"created_at,nullzero,default:now()"`
	UpdatedAt time.Time   `bun:"updated_at,nullzero,default:now()"`
	Email     string      `bun:"email,notnull"`
	Role      domain.Role `bun:"role,notnull"`
	TokenHash string      `bun:"token_hash,notnull"`
	ExpiresAt time.Time   `bun:"expires_at,notnull"`
	UsedAt    *time.Time  `bun:"used_at"`
	InvitedBy int64       `bun:"invited_by,notnull"`
}
