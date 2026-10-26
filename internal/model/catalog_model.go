package model

import (
	"time"

	"github.com/runforyou-ai/ai-gateway/internal/domain"
	"github.com/uptrace/bun"
)

// CatalogModel 对应对外模型目录。
type CatalogModel struct {
	bun.BaseModel `bun:"table:models,alias:m"`

	ID          int64            `bun:"id,pk,autoincrement"`
	CreatedAt   time.Time        `bun:"created_at,nullzero,default:now()"`
	UpdatedAt   time.Time        `bun:"updated_at,nullzero,default:now()"`
	Name        string           `bun:"name,notnull"`
	DisplayName string           `bun:"display_name,notnull"`
	Type        domain.ModelType `bun:"type,notnull"`
	Disabled    bool             `bun:"disabled,notnull,default:false"`
}
