package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/runforyou-ai/ai-gateway/internal/auth"
	"github.com/runforyou-ai/ai-gateway/internal/model"
	"github.com/uptrace/bun"
)

var (
	// ErrInvalidResetToken 表示重置令牌无效或已过期。
	ErrInvalidResetToken = errors.New("invalid reset token")
)

// ResetPassword 使用重置令牌设置新密码，并作废该用户全部会话与未用令牌。
func ResetPassword(ctx context.Context, db bun.IDB, rawToken, newPassword string) error {
	token := strings.TrimSpace(rawToken)
	if token == "" {
		return ErrInvalidResetToken
	}
	if len(newPassword) < 8 {
		return ErrWeakPassword
	}

	hash := auth.HashToken(token)
	now := time.Now().UTC()

	return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		var row model.PasswordResetToken
		err := tx.NewSelect().Model(&row).
			Where("token_hash = ?", hash).
			Where("used_at IS NULL").
			Where("expires_at > ?", now).
			Scan(ctx)
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInvalidResetToken
		}
		if err != nil {
			return err
		}

		pwHash, err := auth.HashPassword(newPassword)
		if err != nil {
			return fmt.Errorf("hash password: %w", err)
		}

		res, err := tx.NewUpdate().Model((*model.User)(nil)).
			Set("password = ?", pwHash).
			Where("id = ?", row.UserID).
			Exec(ctx)
		if err != nil {
			return err
		}
		n, _ := res.RowsAffected()
		if n == 0 {
			return ErrUserNotFound
		}

		used := now
		row.UsedAt = &used
		if _, err := tx.NewUpdate().Model(&row).Column("used_at").WherePK().Exec(ctx); err != nil {
			return err
		}

		if _, err := tx.NewUpdate().Model((*model.PasswordResetToken)(nil)).
			Set("used_at = ?", now).
			Where("user_id = ?", row.UserID).
			Where("used_at IS NULL").
			Where("id <> ?", row.ID).
			Exec(ctx); err != nil {
			return err
		}

		if _, err := tx.NewDelete().Model((*model.Session)(nil)).
			Where("user_id = ?", row.UserID).
			Exec(ctx); err != nil {
			return err
		}
		return nil
	})
}
