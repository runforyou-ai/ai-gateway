// Package sysdoc 提供 settings / state 按名文档读写。
package sysdoc

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/runforyou-ai/lite-ai-gateway/internal/model"
	"github.com/uptrace/bun"
)

var (
	// ErrSettingsMissing 表示指定设置组不存在。
	ErrSettingsMissing = errors.New("settings group missing")
	// ErrStateMissing 表示指定状态组不存在。
	ErrStateMissing = errors.New("state group missing")
)

// GetSettings 读取设置组 JSON（不含 blob）。
func GetSettings[T any](ctx context.Context, db bun.IDB, name string, missingOK bool) (value T, ok bool, err error) {
	row := &model.Setting{}
	err = db.NewSelect().Model(row).
		Column("name", "payload", "version", "created_at", "updated_at").
		Where("name = ?", name).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		if missingOK {
			return value, false, nil
		}
		return value, false, ErrSettingsMissing
	}
	if err != nil {
		return value, false, fmt.Errorf("read settings %s: %w", name, err)
	}
	if len(row.Payload) > 0 {
		if err := json.Unmarshal(row.Payload, &value); err != nil {
			return value, false, fmt.Errorf("decode settings %s: %w", name, err)
		}
	}
	return value, true, nil
}

// PutSettings 整组写入设置；blob 为 nil 或空切片时清空二进制列。
func PutSettings[T any](ctx context.Context, db bun.IDB, name string, value T, blob []byte) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("encode settings %s: %w", name, err)
	}
	if len(blob) == 0 {
		blob = nil
	}
	_, err = db.NewInsert().Model(&model.Setting{
		Name: name, Payload: payload, Blob: blob, Version: 1,
	}).
		On("CONFLICT (name) DO UPDATE").
		Set("payload = EXCLUDED.payload").
		Set("blob = EXCLUDED.blob").
		Set("version = settings.version + 1").
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("write settings %s: %w", name, err)
	}
	return nil
}

// GetState 读取状态组。
func GetState[T any](ctx context.Context, db bun.IDB, name string, missingOK bool) (value T, ok bool, err error) {
	row := &model.State{}
	err = db.NewSelect().Model(row).Where("name = ?", name).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		if missingOK {
			return value, false, nil
		}
		return value, false, ErrStateMissing
	}
	if err != nil {
		return value, false, fmt.Errorf("read state %s: %w", name, err)
	}
	if len(row.Payload) > 0 {
		if err := json.Unmarshal(row.Payload, &value); err != nil {
			return value, false, fmt.Errorf("decode state %s: %w", name, err)
		}
	}
	return value, true, nil
}

// PutState 整组写入状态。
func PutState[T any](ctx context.Context, db bun.IDB, name string, value T) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("encode state %s: %w", name, err)
	}
	_, err = db.NewInsert().Model(&model.State{
		Name: name, Payload: payload, Version: 1,
	}).
		On("CONFLICT (name) DO UPDATE").
		Set("payload = EXCLUDED.payload").
		Set("version = state.version + 1").
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("write state %s: %w", name, err)
	}
	return nil
}

// StateExists 判断状态组是否存在。
func StateExists(ctx context.Context, db bun.IDB, name string) (bool, error) {
	return db.NewSelect().Model((*model.State)(nil)).Where("name = ?", name).Exists(ctx)
}
