package auth

import (
	"context"

	"github.com/runforyou-ai/ai-gateway/internal/auth"
	"github.com/runforyou-ai/ai-gateway/internal/model"
	"github.com/uptrace/bun"
)

// Logout 按 Cookie 删除会话；无会话时视为成功。
func Logout(ctx context.Context, db *bun.DB, cookieHeader string) error {
	token := auth.TokenFromCookieHeader(cookieHeader)
	if token == "" {
		return nil
	}
	hash := auth.HashToken(token)
	_, err := db.NewDelete().Model((*model.Session)(nil)).Where("token_hash = ?", hash).Exec(ctx)
	return err
}
