package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/runforyou-ai/ai-gateway/internal/auth"
	"github.com/runforyou-ai/ai-gateway/internal/model"
	"github.com/runforyou-ai/support/str"
	"github.com/uptrace/bun"
)

var (
	// ErrWrongPassword 表示当前密码不正确。
	ErrWrongPassword = errors.New("wrong password")
	// ErrUserNotFound 表示用户不存在。
	ErrUserNotFound = errors.New("user not found")
)

// UpdateProfileInput 是更新个人资料参数。
type UpdateProfileInput struct {
	Name            string
	Email           string
	CurrentPassword string
	NewPassword     string
}

// UpdateProfile 更新当前用户的姓名、邮箱；若填写新密码则校验当前密码后改密。
func UpdateProfile(ctx context.Context, db bun.IDB, userID int64, in UpdateProfileInput) (UserInfo, error) {
	name := strings.TrimSpace(in.Name)
	email := strings.TrimSpace(strings.ToLower(in.Email))
	if name == "" || len([]rune(name)) > 64 {
		return UserInfo{}, ErrInvalidName
	}
	if !str.IsEmail(email) {
		return UserInfo{}, ErrInvalidEmail
	}

	var user model.User
	err := db.NewSelect().Model(&user).Where("id = ?", userID).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return UserInfo{}, ErrUserNotFound
	}
	if err != nil {
		return UserInfo{}, err
	}

	if email != user.Email {
		var existing model.User
		err := db.NewSelect().Model(&existing).Where("email = ?", email).Scan(ctx)
		if err == nil {
			return UserInfo{}, ErrEmailTaken
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return UserInfo{}, err
		}
	}

	user.Name = name
	user.Email = email
	cols := []string{"name", "email"}

	if in.NewPassword != "" {
		if len(in.NewPassword) < 8 {
			return UserInfo{}, ErrWeakPassword
		}
		if !auth.CheckPassword(user.Password, in.CurrentPassword) {
			return UserInfo{}, ErrWrongPassword
		}
		hash, err := auth.HashPassword(in.NewPassword)
		if err != nil {
			return UserInfo{}, fmt.Errorf("hash password: %w", err)
		}
		user.Password = hash
		cols = append(cols, "password")
	}

	if _, err := db.NewUpdate().Model(&user).Column(cols...).WherePK().Exec(ctx); err != nil {
		return UserInfo{}, err
	}

	return UserInfo{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}
