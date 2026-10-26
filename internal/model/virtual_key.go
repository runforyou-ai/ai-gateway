package model

import (
	"time"

	"github.com/uptrace/bun"
)

// VirtualKey 对应调用 /v1 的虚拟密钥。
type VirtualKey struct {
	bun.BaseModel `bun:"table:virtual_keys,alias:vk"`

	ID          int64      `bun:"id,pk,autoincrement"`
	CreatedAt   time.Time  `bun:"created_at,nullzero,default:now()"`
	UpdatedAt   time.Time  `bun:"updated_at,nullzero,default:now()"`
	UserID      int64      `bun:"user_id,notnull"`
	Name        string     `bun:"name,notnull"`
	Prefix      string     `bun:"prefix,notnull"`
	TokenHash   string     `bun:"token_hash,notnull"`
	Disabled    bool       `bun:"disabled,notnull,default:false"`
	MaxBudget   *float64   `bun:"max_budget"`
	Spent       float64    `bun:"spent,notnull,default:0"`
	BudgetReset string     `bun:"budget_reset,notnull,default:''"`
	ExpiresIn   string     `bun:"expires_in,notnull,default:''"`
	ExpiresAt   *time.Time `bun:"expires_at"`
	LastUsedAt  *time.Time `bun:"last_used_at"`
}
