package model

import (
	"encoding/json"
	"time"

	"github.com/uptrace/bun"
)

// Setting 对应 settings 表的一组文档。
type Setting struct {
	bun.BaseModel `bun:"table:settings,alias:setg"`

	Name      string          `bun:"name,pk"`
	Payload   json.RawMessage `bun:"payload,type:jsonb"`
	Blob      []byte          `bun:"blob"`
	Version   int64           `bun:"version"`
	CreatedAt time.Time       `bun:"created_at,nullzero,default:now()"`
	UpdatedAt time.Time       `bun:"updated_at,nullzero,default:now()"`
}
