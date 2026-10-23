package auth

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/runforyou-ai/lite-ai-gateway/internal/auth"
	"github.com/runforyou-ai/lite-ai-gateway/internal/domain"
	"github.com/runforyou-ai/lite-ai-gateway/internal/model"
	"github.com/uptrace/bun"
)

var (
	// ErrInvalidCredentials 表示邮箱或密码错误。
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// LoginInput 是登录参数。
type LoginInput struct {
	Email    string
	Password string
	Remember bool
}

// LoginResult 含明文令牌（仅用于 Set-Cookie）与用户信息。
type LoginResult struct {
	Token    string
	Remember bool
	User     UserInfo
}

// UserInfo 是对外暴露的用户信息。
type UserInfo struct {
	ID    int64
	Name  string
	Email string
	Role  domain.Role
}

// Login 校验密码并创建会话，返回明文令牌。
func Login(ctx context.Context, db *bun.DB, in LoginInput) (LoginResult, error) {
	email := strings.TrimSpace(strings.ToLower(in.Email))
	if email == "" || in.Password == "" {
		return LoginResult{}, ErrInvalidCredentials
	}
	var user model.User
	err := db.NewSelect().Model(&user).Where("email = ?", email).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return LoginResult{}, ErrInvalidCredentials
	}
	if err != nil {
		return LoginResult{}, err
	}
	if !auth.CheckPassword(user.Password, in.Password) {
		return LoginResult{}, ErrInvalidCredentials
	}

	token, err := auth.NewSessionToken()
	if err != nil {
		return LoginResult{}, err
	}
	ttl := auth.SessionTTLFor(in.Remember)
	sess := &model.Session{
		UserID:    user.ID,
		TokenHash: auth.HashToken(token),
		ExpiresAt: time.Now().Add(ttl),
	}
	if _, err := db.NewInsert().Model(sess).Exec(ctx); err != nil {
		return LoginResult{}, err
	}
	return LoginResult{
		Token:    token,
		Remember: in.Remember,
		User: UserInfo{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}
