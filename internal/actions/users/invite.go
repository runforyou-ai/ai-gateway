package users

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
	"github.com/runforyou-ai/ai-gateway/internal/domain"
	"github.com/runforyou-ai/ai-gateway/internal/email"
	"github.com/runforyou-ai/ai-gateway/internal/model"
	"github.com/runforyou-ai/support/str"
	"github.com/uptrace/bun"
)

const inviteTTL = 7 * 24 * time.Hour

var (
	// ErrEmailNotConfigured 表示未配置发信。
	ErrEmailNotConfigured = errors.New("email not configured")
	// ErrInvalidInvite 表示邀请无效或已过期。
	ErrInvalidInvite = errors.New("invalid invite")
)

// InviteInput 是创建邀请参数。
type InviteInput struct {
	Email string
	Role  domain.Role
}

// InviteLookup 是公开查询邀请的结果。
type InviteLookup struct {
	Email string
	Role  domain.Role
}

// AcceptInviteInput 是接受邀请参数。
type AcceptInviteInput struct {
	Token    string
	Name     string
	Password string
}

// CreateInvite 创建邀请并发送邮件；同邮箱未使用的旧邀请会被作废。
func CreateInvite(ctx context.Context, db bun.IDB, mailer email.Sender, publicURL string, actorID int64, in InviteInput) error {
	emailAddr := strings.TrimSpace(strings.ToLower(in.Email))
	if !str.IsEmail(emailAddr) {
		return ErrInvalidEmail
	}
	if !in.Role.Valid() {
		return ErrInvalidRole
	}
	if mailer == nil {
		return ErrEmailNotConfigured
	}
	base := strings.TrimRight(strings.TrimSpace(publicURL), "/")
	if base == "" {
		return fmt.Errorf("public url is required")
	}
	if err := ensureEmailFree(ctx, db, emailAddr, 0); err != nil {
		return err
	}

	token, hash := auth.NewSessionToken()
	row := &model.UserInvite{
		Email:     emailAddr,
		Role:      in.Role,
		TokenHash: hash,
		ExpiresAt: time.Now().UTC().Add(inviteTTL),
		InvitedBy: actorID,
	}
	err := db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewDelete().Model((*model.UserInvite)(nil)).
			Where("lower(email) = ?", emailAddr).
			Where("used_at IS NULL").
			Exec(ctx); err != nil {
			return err
		}
		_, err := tx.NewInsert().Model(row).Exec(ctx)
		return err
	})
	if err != nil {
		return err
	}

	link := base + "/accept-invite?token=" + url.QueryEscape(token)
	safeLink := html.EscapeString(link)
	err = mailer.Send(ctx, email.Message{
		To:      emailAddr,
		Subject: "You are invited to AI Gateway",
		Text:    "You have been invited. Open this link to set your password (valid for 7 days):\n\n" + link + "\n",
		HTML:    "<p>You have been invited. Open this link to set your password (valid for 7 days):</p><p><a href=\"" + safeLink + "\">" + safeLink + "</a></p>",
	})
	if err != nil {
		_, _ = db.NewDelete().Model(row).WherePK().Exec(ctx)
		if errors.Is(err, email.ErrNotConfigured) {
			return ErrEmailNotConfigured
		}
		return fmt.Errorf("send invite email: %w", err)
	}
	return nil
}

// LookupInvite 按明文令牌查询未使用且未过期的邀请。
func LookupInvite(ctx context.Context, db bun.IDB, rawToken string) (InviteLookup, error) {
	row, err := findValidInvite(ctx, db, rawToken)
	if err != nil {
		return InviteLookup{}, err
	}
	return InviteLookup{Email: row.Email, Role: row.Role}, nil
}

// AcceptInvite 接受邀请并创建账号。
func AcceptInvite(ctx context.Context, db bun.IDB, in AcceptInviteInput) (User, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" || len([]rune(name)) > 64 {
		return User{}, ErrInvalidName
	}
	if len(in.Password) < 8 {
		return User{}, ErrWeakPassword
	}

	hash, err := auth.HashPassword(in.Password)
	if err != nil {
		return User{}, fmt.Errorf("hash password: %w", err)
	}

	var created User
	err = db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		row, err := findValidInvite(ctx, tx, in.Token)
		if err != nil {
			return err
		}
		if err := ensureEmailFree(ctx, tx, row.Email, 0); err != nil {
			return err
		}

		user := &model.User{
			Name:     name,
			Email:    row.Email,
			Password: hash,
			Role:     row.Role,
			Disabled: false,
		}
		if _, err := tx.NewInsert().Model(user).Exec(ctx); err != nil {
			return err
		}
		now := time.Now().UTC()
		res, err := tx.NewUpdate().Model((*model.UserInvite)(nil)).
			Set("used_at = ?", now).
			Where("id = ?", row.ID).
			Where("used_at IS NULL").
			Exec(ctx)
		if err != nil {
			return err
		}
		n, _ := res.RowsAffected()
		if n == 0 {
			return ErrInvalidInvite
		}
		created = toUser(*user)
		return nil
	})
	if err != nil {
		return User{}, err
	}
	return created, nil
}

func findValidInvite(ctx context.Context, db bun.IDB, rawToken string) (model.UserInvite, error) {
	token := strings.TrimSpace(rawToken)
	if token == "" {
		return model.UserInvite{}, ErrInvalidInvite
	}
	hash := auth.HashToken(token)
	now := time.Now().UTC()
	var row model.UserInvite
	err := db.NewSelect().Model(&row).
		Where("token_hash = ?", hash).
		Where("used_at IS NULL").
		Where("expires_at > ?", now).
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return model.UserInvite{}, ErrInvalidInvite
	}
	if err != nil {
		return model.UserInvite{}, err
	}
	return row, nil
}
