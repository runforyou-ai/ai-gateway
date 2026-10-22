package model

import (
	"encoding/json"
	"time"

	"github.com/uptrace/bun"
)

// State 对应 state 表的一组文档。
type State struct {
	bun.BaseModel `bun:"table:state,alias:st"`

	Name      string          `bun:"name,pk"`
	Payload   json.RawMessage `bun:"payload,type:jsonb"`
	Version   int64           `bun:"version"`
	CreatedAt time.Time       `bun:"created_at,nullzero,default:now()"`
	UpdatedAt time.Time       `bun:"updated_at,nullzero,default:now()"`
}
