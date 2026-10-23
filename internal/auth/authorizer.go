package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/runforyou-ai/lite-ai-gateway/internal/domain"
	"github.com/runforyou-ai/lite-ai-gateway/internal/model"
	"github.com/uptrace/bun"
)

// SessionAuthorizer 用 Cookie 会话解析 Principal。
type SessionAuthorizer struct {
	DB *bun.DB
}

// Authenticate 校验 Cookie 会话并返回主体；成功时滑动续期。
func (a *SessionAuthorizer) Authenticate(ctx context.Context, cookieHeader string) (Principal, error) {
	token := TokenFromCookieHeader(cookieHeader)
	if token == "" {
		return Principal{}, errors.New("missing session")
	}
	hash := HashToken(token)
	var sess model.Session
	err := a.DB.NewSelect().Model(&sess).Where("token_hash = ?", hash).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return Principal{}, errors.New("invalid session")
	}
	if err != nil {
		return Principal{}, err
	}
	if time.Now().After(sess.ExpiresAt) {
		_, _ = a.DB.NewDelete().Model(&sess).WherePK().Exec(ctx)
		return Principal{}, errors.New("session expired")
	}

	var user model.User
	if err := a.DB.NewSelect().Model(&user).Where("id = ?", sess.UserID).Scan(ctx); err != nil {
		return Principal{}, err
	}
	if !user.Role.Valid() {
		return Principal{}, errors.New("invalid role")
	}

	newExp := time.Now().Add(SessionTTL)
	_, _ = a.DB.NewUpdate().Model(&sess).
		Set("expires_at = ?", newExp).
		WherePK().
		Exec(ctx)

	return Principal{
		UserID:    user.ID,
		SessionID: sess.ID,
		Role:      domain.Role(user.Role),
	}, nil
}
