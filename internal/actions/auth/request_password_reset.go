package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"html"
	"net/url"
	"strings"
	"time"

	"github.com/runforyou-ai/ai-gateway/internal/auth"
	"github.com/runforyou-ai/ai-gateway/internal/email"
	"github.com/runforyou-ai/ai-gateway/internal/model"
	"github.com/runforyou-ai/support/str"
	"github.com/uptrace/bun"
)

const passwordResetTTL = time.Hour

var (
	// ErrEmailNotConfigured 表示未配置发信，无法发送重置邮件。
	ErrEmailNotConfigured = errors.New("email not configured")
)

// RequestPasswordReset 为已存在用户创建重置令牌并发送邮件。
// 邮箱不存在时不报错（防枚举）；格式非法返回 ErrInvalidEmail。
func RequestPasswordReset(ctx context.Context, db bun.IDB, mailer email.Sender, publicURL, rawEmail string) error {
	addr := strings.TrimSpace(strings.ToLower(rawEmail))
	if !str.IsEmail(addr) {
		return ErrInvalidEmail
	}
	if mailer == nil {
		return ErrEmailNotConfigured
	}

	var user model.User
	err := db.NewSelect().Model(&user).Where("email = ?", addr).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}

	base := strings.TrimRight(strings.TrimSpace(publicURL), "/")
	if base == "" {
		return fmt.Errorf("public url is required")
	}

	token, hash := auth.NewSessionToken()
	row := &model.PasswordResetToken{
		UserID:    user.ID,
		TokenHash: hash,
		ExpiresAt: time.Now().UTC().Add(passwordResetTTL),
	}
	if _, err := db.NewInsert().Model(row).Exec(ctx); err != nil {
		return err
	}

	link := base + "/reset-password?token=" + url.QueryEscape(token)

	safeLink := html.EscapeString(link)
	err = mailer.Send(ctx, email.Message{
		To:      user.Email,
		Subject: "Reset your password",
		Text:    "Use this link to reset your password (valid for 1 hour):\n\n" + link + "\n",
		HTML:    "<p>Use this link to reset your password (valid for 1 hour):</p><p><a href=\"" + safeLink + "\">" + safeLink + "</a></p>",
	})
	if err != nil {
		_, _ = db.NewDelete().Model(row).WherePK().Exec(ctx)
		if errors.Is(err, email.ErrNotConfigured) {
			return ErrEmailNotConfigured
		}
		return fmt.Errorf("send reset email: %w", err)
	}
	return nil
}
