package models

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/runforyou-ai/ai-gateway/internal/domain"
	"github.com/runforyou-ai/ai-gateway/internal/model"
	"github.com/uptrace/bun"
)

var (
	// ErrInvalidName 表示对外模型 id 无效。
	ErrInvalidName = errors.New("invalid name")
	// ErrInvalidDisplayName 表示展示名称无效。
	ErrInvalidDisplayName = errors.New("invalid display name")
	// ErrInvalidType 表示模型类型无效。
	ErrInvalidType = errors.New("invalid type")
	// ErrNameTaken 表示对外模型 id 已被占用。
	ErrNameTaken = errors.New("name taken")
	// ErrNotFound 表示模型不存在。
	ErrNotFound = errors.New("model not found")
)

// Model 是对外目录模型视图。
type Model struct {
	ID          int64
	Name        string
	DisplayName string
	Type        domain.ModelType
	Disabled    bool
	CreatedAt   time.Time
}

// SupportedProtocols 返回该模型对外支持的协议全集（由类型推导）。
func (m Model) SupportedProtocols() []domain.Protocol {
	return m.Type.Protocols()
}

// CreateInput 是创建参数。
type CreateInput struct {
	Name        string
	DisplayName string
	Type        domain.ModelType
	Disabled    bool
}

// UpdateInput 是更新参数。
type UpdateInput struct {
	Name        string
	DisplayName string
	Type        domain.ModelType
	Disabled    bool
}

// List 分页列出目录模型。
func List(ctx context.Context, db bun.IDB, start, end int) ([]Model, int64, error) {
	if start < 0 {
		start = 0
	}
	if end <= start {
		end = start + 20
	}
	limit := end - start
	if limit > 100 {
		limit = 100
	}

	total, err := db.NewSelect().Model((*model.CatalogModel)(nil)).Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	var rows []model.CatalogModel
	err = db.NewSelect().Model(&rows).
		OrderExpr("id DESC").
		Offset(int64(start)).
		Limit(int64(limit)).
		Scan(ctx)
	if err != nil {
		return nil, 0, err
	}
	out := make([]Model, 0, len(rows))
	for _, row := range rows {
		out = append(out, toModel(row))
	}
	return out, int64(total), nil
}

// Get 按 id 获取目录模型。
func Get(ctx context.Context, db bun.IDB, id int64) (Model, error) {
	row, err := getByID(ctx, db, id)
	if err != nil {
		return Model{}, err
	}
	return toModel(row), nil
}

// Create 创建目录模型。
func Create(ctx context.Context, db bun.IDB, in CreateInput) (Model, error) {
	name, displayName, typ, err := normalizeIdentity(in.Name, in.DisplayName, in.Type)
	if err != nil {
		return Model{}, err
	}
	taken, err := nameExists(ctx, db, name, 0)
	if err != nil {
		return Model{}, err
	}
	if taken {
		return Model{}, ErrNameTaken
	}
	row := &model.CatalogModel{
		Name:        name,
		DisplayName: displayName,
		Type:        typ,
		Disabled:    in.Disabled,
	}
	if _, err := db.NewInsert().Model(row).Exec(ctx); err != nil {
		return Model{}, err
	}
	return toModel(*row), nil
}

// Update 更新目录模型。
func Update(ctx context.Context, db bun.IDB, id int64, in UpdateInput) (Model, error) {
	name, displayName, typ, err := normalizeIdentity(in.Name, in.DisplayName, in.Type)
	if err != nil {
		return Model{}, err
	}
	row, err := getByID(ctx, db, id)
	if err != nil {
		return Model{}, err
	}
	if name != row.Name {
		taken, err := nameExists(ctx, db, name, id)
		if err != nil {
			return Model{}, err
		}
		if taken {
			return Model{}, ErrNameTaken
		}
	}
	row.Name = name
	row.DisplayName = displayName
	row.Type = typ
	row.Disabled = in.Disabled
	if _, err := db.NewUpdate().Model(&row).
		Column("name", "display_name", "type", "disabled").
		WherePK().
		Exec(ctx); err != nil {
		return Model{}, err
	}
	return toModel(row), nil
}

// Delete 删除目录模型。
func Delete(ctx context.Context, db bun.IDB, id int64) error {
	res, err := db.NewDelete().Model((*model.CatalogModel)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func getByID(ctx context.Context, db bun.IDB, id int64) (model.CatalogModel, error) {
	var row model.CatalogModel
	err := db.NewSelect().Model(&row).Where("id = ?", id).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return model.CatalogModel{}, ErrNotFound
	}
	if err != nil {
		return model.CatalogModel{}, err
	}
	return row, nil
}

func nameExists(ctx context.Context, db bun.IDB, name string, excludeID int64) (bool, error) {
	q := db.NewSelect().Model((*model.CatalogModel)(nil)).Where("name = ?", name)
	if excludeID > 0 {
		q = q.Where("id <> ?", excludeID)
	}
	return q.Exists(ctx)
}

func normalizeIdentity(name, displayName string, typ domain.ModelType) (string, string, domain.ModelType, error) {
	name = strings.TrimSpace(name)
	displayName = strings.TrimSpace(displayName)
	if !validModelName(name) {
		return "", "", "", ErrInvalidName
	}
	if displayName == "" || utf8.RuneCountInString(displayName) > 128 {
		return "", "", "", ErrInvalidDisplayName
	}
	if !typ.Valid() {
		return "", "", "", ErrInvalidType
	}
	return name, displayName, typ, nil
}

func validModelName(name string) bool {
	if name == "" || utf8.RuneCountInString(name) > 128 {
		return false
	}
	for i, r := range name {
		switch {
		case unicode.IsLetter(r), unicode.IsDigit(r):
			continue
		case r == '.' || r == '_' || r == '-' || r == '/' || r == ':':
			if i == 0 {
				return false
			}
			continue
		default:
			return false
		}
	}
	return true
}

func toModel(row model.CatalogModel) Model {
	return Model{
		ID:          row.ID,
		Name:        row.Name,
		DisplayName: row.DisplayName,
		Type:        row.Type,
		Disabled:    row.Disabled,
		CreatedAt:   row.CreatedAt,
	}
}
