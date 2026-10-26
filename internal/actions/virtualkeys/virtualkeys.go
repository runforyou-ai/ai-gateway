package virtualkeys

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/runforyou-ai/ai-gateway/internal/auth"
	"github.com/runforyou-ai/ai-gateway/internal/model"
	"github.com/uptrace/bun"
)

const (
	keyPrefixTag = "vk_"
	keyPrefixLen = 11 // "vk_" + 8
)

var (
	// ErrInvalidName 表示名称无效。
	ErrInvalidName = errors.New("invalid name")
	// ErrInvalidMaxBudget 表示消费上限无效。
	ErrInvalidMaxBudget = errors.New("invalid max budget")
	// ErrInvalidExpiresIn 表示过期时长无效。
	ErrInvalidExpiresIn = errors.New("invalid expires in")
	// ErrInvalidBudgetReset 表示额度重置周期无效。
	ErrInvalidBudgetReset = errors.New("invalid budget reset")
	// ErrNotFound 表示虚拟密钥不存在或不属于当前用户。
	ErrNotFound = errors.New("virtual key not found")
)

// Key 是对外视图（不含明文）。
type Key struct {
	ID          int64
	Name        string
	Prefix      string
	Disabled    bool
	MaxBudget   *float64
	Spent       float64
	BudgetReset string
	ExpiresIn   string
	ExpiresAt   *time.Time
	LastUsedAt  *time.Time
	CreatedAt   time.Time
}

// Created 是创建结果，含一次性明文。
type Created struct {
	Key
	Secret string
}

// CreateInput 是创建参数。
type CreateInput struct {
	Name        string
	MaxBudget   string
	ExpiresIn   string
	BudgetReset string
}

// UpdateInput 是更新参数。
type UpdateInput struct {
	Name        string
	Disabled    bool
	MaxBudget   string
	ExpiresIn   string
	BudgetReset string
}

// List 列出用户自己的虚拟密钥。
func List(ctx context.Context, db bun.IDB, userID int64, start, end int) ([]Key, int64, error) {
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

	total, err := db.NewSelect().Model((*model.VirtualKey)(nil)).Where("user_id = ?", userID).Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	var rows []model.VirtualKey
	err = db.NewSelect().Model(&rows).
		Where("user_id = ?", userID).
		OrderExpr("id DESC").
		Offset(int64(start)).
		Limit(int64(limit)).
		Scan(ctx)
	if err != nil {
		return nil, 0, err
	}
	out := make([]Key, 0, len(rows))
	for _, row := range rows {
		out = append(out, toKey(row))
	}
	return out, int64(total), nil
}

// Get 获取属于 userID 的虚拟密钥。
func Get(ctx context.Context, db bun.IDB, userID, id int64) (Key, error) {
	row, err := getOwned(ctx, db, userID, id)
	if err != nil {
		return Key{}, err
	}
	return toKey(row), nil
}

// Create 创建虚拟密钥；明文只在返回值中给出一次。
func Create(ctx context.Context, db bun.IDB, userID int64, in CreateInput) (Created, error) {
	name, err := normalizeName(in.Name)
	if err != nil {
		return Created{}, err
	}
	maxBudget, err := parseMaxBudget(in.MaxBudget)
	if err != nil {
		return Created{}, err
	}
	expiresIn, expiresAt, err := parseExpiresIn(in.ExpiresIn, time.Now().UTC())
	if err != nil {
		return Created{}, err
	}
	budgetReset, err := parseBudgetReset(in.BudgetReset)
	if err != nil {
		return Created{}, err
	}
	token, _ := auth.NewSessionToken()
	secret := keyPrefixTag + token
	prefix := secret
	if utf8.RuneCountInString(prefix) > keyPrefixLen {
		prefix = string([]rune(prefix)[:keyPrefixLen])
	}
	row := &model.VirtualKey{
		UserID:      userID,
		Name:        name,
		Prefix:      prefix,
		TokenHash:   auth.HashToken(secret),
		Disabled:    false,
		MaxBudget:   maxBudget,
		BudgetReset: budgetReset,
		ExpiresIn:   expiresIn,
		ExpiresAt:   expiresAt,
	}
	if _, err := db.NewInsert().Model(row).Exec(ctx); err != nil {
		return Created{}, err
	}
	return Created{Key: toKey(*row), Secret: secret}, nil
}

// Update 更新属于 userID 的虚拟密钥。
func Update(ctx context.Context, db bun.IDB, userID, id int64, in UpdateInput) (Key, error) {
	name, err := normalizeName(in.Name)
	if err != nil {
		return Key{}, err
	}
	maxBudget, err := parseMaxBudget(in.MaxBudget)
	if err != nil {
		return Key{}, err
	}
	budgetReset, err := parseBudgetReset(in.BudgetReset)
	if err != nil {
		return Key{}, err
	}
	row, err := getOwned(ctx, db, userID, id)
	if err != nil {
		return Key{}, err
	}
	expiresIn := strings.TrimSpace(in.ExpiresIn)
	if expiresIn != row.ExpiresIn {
		code, at, err := parseExpiresIn(expiresIn, time.Now().UTC())
		if err != nil {
			return Key{}, err
		}
		row.ExpiresIn = code
		row.ExpiresAt = at
	}
	row.Name = name
	row.Disabled = in.Disabled
	row.MaxBudget = maxBudget
	row.BudgetReset = budgetReset
	if _, err := db.NewUpdate().Model(&row).
		Column("name", "disabled", "max_budget", "budget_reset", "expires_in", "expires_at").
		WherePK().
		Exec(ctx); err != nil {
		return Key{}, err
	}
	return toKey(row), nil
}

// Delete 删除属于 userID 的虚拟密钥。
func Delete(ctx context.Context, db bun.IDB, userID, id int64) error {
	res, err := db.NewDelete().Model((*model.VirtualKey)(nil)).
		Where("id = ?", id).
		Where("user_id = ?", userID).
		Exec(ctx)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func getOwned(ctx context.Context, db bun.IDB, userID, id int64) (model.VirtualKey, error) {
	var row model.VirtualKey
	err := db.NewSelect().Model(&row).
		Where("id = ?", id).
		Where("user_id = ?", userID).
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return model.VirtualKey{}, ErrNotFound
	}
	if err != nil {
		return model.VirtualKey{}, err
	}
	return row, nil
}

func normalizeName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" || utf8.RuneCountInString(name) > 64 {
		return "", ErrInvalidName
	}
	return name, nil
}

func parseMaxBudget(raw string) (*float64, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil || v <= 0 {
		return nil, ErrInvalidMaxBudget
	}
	return &v, nil
}

func parseExpiresIn(raw string, now time.Time) (string, *time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", nil, nil
	}
	var d time.Duration
	switch raw {
	case "1h":
		d = time.Hour
	case "1d":
		d = 24 * time.Hour
	case "3d":
		d = 3 * 24 * time.Hour
	case "7d":
		d = 7 * 24 * time.Hour
	case "30d":
		d = 30 * 24 * time.Hour
	case "90d":
		d = 90 * 24 * time.Hour
	case "180d":
		d = 180 * 24 * time.Hour
	case "365d":
		d = 365 * 24 * time.Hour
	default:
		return "", nil, ErrInvalidExpiresIn
	}
	at := now.Add(d)
	return raw, &at, nil
}

func parseBudgetReset(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	switch raw {
	case "", "daily", "weekly", "monthly":
		return raw, nil
	default:
		return "", ErrInvalidBudgetReset
	}
}

func toKey(row model.VirtualKey) Key {
	return Key{
		ID:          row.ID,
		Name:        row.Name,
		Prefix:      row.Prefix,
		Disabled:    row.Disabled,
		MaxBudget:   row.MaxBudget,
		Spent:       row.Spent,
		BudgetReset: row.BudgetReset,
		ExpiresIn:   row.ExpiresIn,
		ExpiresAt:   row.ExpiresAt,
		LastUsedAt:  row.LastUsedAt,
		CreatedAt:   row.CreatedAt,
	}
}
