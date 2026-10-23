package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/runforyou-ai/lite-ai-gateway/internal/auth"
	"github.com/runforyou-ai/lite-ai-gateway/internal/domain"
	"github.com/runforyou-ai/lite-ai-gateway/internal/model"
	"github.com/runforyou-ai/lite-ai-gateway/internal/sysdoc"
	"github.com/uptrace/bun"
)

var (
	// ErrEmailTaken 表示邮箱已被占用。
	ErrEmailTaken = errors.New("email taken")
	// ErrRegistrationClosed 表示未开放注册（已有用户且未开启注册）。
	ErrRegistrationClosed = errors.New("registration closed")
	// ErrWeakPassword 表示密码不符合要求。
	ErrWeakPassword = errors.New("weak password")
	// ErrInvalidEmail 表示邮箱无效。
	ErrInvalidEmail = errors.New("invalid email")
)

// RegisterInput 是注册参数。
type RegisterInput struct {
	Email    string
	Password string
}

// RegisterResult 是注册结果。
type RegisterResult struct {
	Role domain.Role
}

// Register 注册账号；第一个用户自动成为 admin，其后为 user（需开放注册）。
func Register(ctx context.Context, db *bun.DB, in RegisterInput) (RegisterResult, error) {
	email := strings.TrimSpace(strings.ToLower(in.Email))
	password := in.Password
	if email == "" || !strings.Contains(email, "@") {
		return RegisterResult{}, ErrInvalidEmail
	}
	if len(password) < 8 {
		return RegisterResult{}, ErrWeakPassword
	}

	var existing model.User
	err := db.NewSelect().Model(&existing).Where("email = ?", email).Scan(ctx)
	if err == nil {
		return RegisterResult{}, ErrEmailTaken
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return RegisterResult{}, err
	}

	n, err := db.NewSelect().Model((*model.User)(nil)).Count(ctx)
	if err != nil {
		return RegisterResult{}, err
	}
	if n > 0 {
		site, ok, err := sysdoc.GetSettings[sysdoc.Site](ctx, db, sysdoc.SettingsSite, true)
		if err != nil {
			return RegisterResult{}, err
		}
		open := false
		if ok {
			open = site.RegistrationOpen
		}
		if !open {
			return RegisterResult{}, ErrRegistrationClosed
		}
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		return RegisterResult{}, fmt.Errorf("hash password: %w", err)
	}

	role := domain.RoleUser
	if n == 0 {
		role = domain.RoleAdmin
	}

	err = db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		user := &model.User{
			Email:    email,
			Password: hash,
			Role:     role,
		}
		if _, err := tx.NewInsert().Model(user).Exec(ctx); err != nil {
			return err
		}
		if n == 0 {
			site := sysdoc.DefaultSite()
			if err := sysdoc.PutSettings(ctx, tx, sysdoc.SettingsSite, site, nil); err != nil {
				return err
			}
			if err := sysdoc.PutState(ctx, tx, sysdoc.StateInstallation, sysdoc.Installation{Installed: true}); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return RegisterResult{}, err
	}
	return RegisterResult{Role: role}, nil
}
