package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/runforyou-ai/ai-gateway/internal/auth"
	"github.com/runforyou-ai/ai-gateway/internal/domain"
	"github.com/runforyou-ai/ai-gateway/internal/model"
	"github.com/runforyou-ai/support/str"
	"github.com/uptrace/bun"
)

var (
	// ErrInvalidName 表示用户名无效。
	ErrInvalidName = errors.New("invalid name")
	// ErrInvalidEmail 表示邮箱无效。
	ErrInvalidEmail = errors.New("invalid email")
	// ErrWeakPassword 表示密码不符合要求。
	ErrWeakPassword = errors.New("weak password")
	// ErrInvalidRole 表示角色无效。
	ErrInvalidRole = errors.New("invalid role")
	// ErrEmailTaken 表示邮箱已被占用。
	ErrEmailTaken = errors.New("email taken")
	// ErrNotFound 表示用户不存在。
	ErrNotFound = errors.New("user not found")
	// ErrCannotDeleteSelf 表示不能删除自己。
	ErrCannotDeleteSelf = errors.New("cannot delete self")
	// ErrCannotDisableSelf 表示不能禁用自己。
	ErrCannotDisableSelf = errors.New("cannot disable self")
	// ErrLastAdmin 表示操作会导致没有管理员。
	ErrLastAdmin = errors.New("last admin")
)

// User 是对外用户视图。
type User struct {
	ID       int64
	Name     string
	Email    string
	Role     domain.Role
	Disabled bool
}

// CreateInput 是创建用户参数。
type CreateInput struct {
	Name     string
	Email    string
	Password string
	Role     domain.Role
	Disabled bool
}

// UpdateInput 是更新用户参数；Password 为空表示不改密。
type UpdateInput struct {
	Name     string
	Email    string
	Password string
	Role     domain.Role
	Disabled bool
}

// List 按 _start/_end 半开区间分页列出用户。
func List(ctx context.Context, db bun.IDB, start, end int) ([]User, int64, error) {
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

	total, err := db.NewSelect().Model((*model.User)(nil)).Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	var rows []model.User
	err = db.NewSelect().Model(&rows).
		OrderExpr("id ASC").
		Offset(int64(start)).
		Limit(int64(limit)).
		Scan(ctx)
	if err != nil {
		return nil, 0, err
	}
	out := make([]User, 0, len(rows))
	for _, row := range rows {
		out = append(out, toUser(row))
	}
	return out, int64(total), nil
}

// Get 按 id 获取用户。
func Get(ctx context.Context, db bun.IDB, id int64) (User, error) {
	var row model.User
	err := db.NewSelect().Model(&row).Where("id = ?", id).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrNotFound
	}
	if err != nil {
		return User{}, err
	}
	return toUser(row), nil
}

// Create 创建用户。
func Create(ctx context.Context, db bun.IDB, in CreateInput) (User, error) {
	name, email, role, err := normalizeIdentity(in.Name, in.Email, in.Role)
	if err != nil {
		return User{}, err
	}
	if len(in.Password) < 8 {
		return User{}, ErrWeakPassword
	}
	if err := ensureEmailFree(ctx, db, email, 0); err != nil {
		return User{}, err
	}
	hash, err := auth.HashPassword(in.Password)
	if err != nil {
		return User{}, fmt.Errorf("hash password: %w", err)
	}
	row := &model.User{
		Name:     name,
		Email:    email,
		Password: hash,
		Role:     role,
		Disabled: in.Disabled,
	}
	if _, err := db.NewInsert().Model(row).Exec(ctx); err != nil {
		return User{}, err
	}
	return toUser(*row), nil
}

// Update 更新用户；actorID 用于禁止禁用自己。
func Update(ctx context.Context, db bun.IDB, actorID, id int64, in UpdateInput) (User, error) {
	name, email, role, err := normalizeIdentity(in.Name, in.Email, in.Role)
	if err != nil {
		return User{}, err
	}
	if in.Password != "" && len(in.Password) < 8 {
		return User{}, ErrWeakPassword
	}

	var row model.User
	err = db.NewSelect().Model(&row).Where("id = ?", id).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrNotFound
	}
	if err != nil {
		return User{}, err
	}

	if in.Disabled && id == actorID {
		return User{}, ErrCannotDisableSelf
	}

	if row.Role == domain.RoleAdmin && (role != domain.RoleAdmin || in.Disabled) {
		ok, err := hasOtherAdmin(ctx, db, id)
		if err != nil {
			return User{}, err
		}
		if !ok {
			return User{}, ErrLastAdmin
		}
	}

	if err := ensureEmailFree(ctx, db, email, id); err != nil {
		return User{}, err
	}

	row.Name = name
	row.Email = email
	row.Role = role
	row.Disabled = in.Disabled
	cols := []string{"name", "email", "role", "disabled"}
	if in.Password != "" {
		hash, err := auth.HashPassword(in.Password)
		if err != nil {
			return User{}, fmt.Errorf("hash password: %w", err)
		}
		row.Password = hash
		cols = append(cols, "password")
	}
	if _, err := db.NewUpdate().Model(&row).Column(cols...).WherePK().Exec(ctx); err != nil {
		return User{}, err
	}
	if in.Disabled {
		_, _ = db.NewDelete().Model((*model.Session)(nil)).Where("user_id = ?", id).Exec(ctx)
	}
	return toUser(row), nil
}

// Delete 删除用户。
func Delete(ctx context.Context, db bun.IDB, actorID, id int64) error {
	if id == actorID {
		return ErrCannotDeleteSelf
	}
	var row model.User
	err := db.NewSelect().Model(&row).Where("id = ?", id).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if row.Role == domain.RoleAdmin {
		ok, err := hasOtherAdmin(ctx, db, id)
		if err != nil {
			return err
		}
		if !ok {
			return ErrLastAdmin
		}
	}
	res, err := db.NewDelete().Model((*model.User)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func normalizeIdentity(name, email string, role domain.Role) (string, string, domain.Role, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(strings.ToLower(email))
	if name == "" || len([]rune(name)) > 64 {
		return "", "", "", ErrInvalidName
	}
	if !str.IsEmail(email) {
		return "", "", "", ErrInvalidEmail
	}
	if !role.Valid() {
		return "", "", "", ErrInvalidRole
	}
	return name, email, role, nil
}

func ensureEmailFree(ctx context.Context, db bun.IDB, email string, exceptID int64) error {
	q := db.NewSelect().Model((*model.User)(nil)).Where("email = ?", email)
	if exceptID > 0 {
		q = q.Where("id <> ?", exceptID)
	}
	n, err := q.Count(ctx)
	if err != nil {
		return err
	}
	if n > 0 {
		return ErrEmailTaken
	}
	return nil
}

func hasOtherAdmin(ctx context.Context, db bun.IDB, exceptID int64) (bool, error) {
	n, err := db.NewSelect().Model((*model.User)(nil)).
		Where("role = ?", domain.RoleAdmin).
		Where("disabled = ?", false).
		Where("id <> ?", exceptID).
		Count(ctx)
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func toUser(row model.User) User {
	return User{
		ID:       row.ID,
		Name:     row.Name,
		Email:    row.Email,
		Role:     row.Role,
		Disabled: row.Disabled,
	}
}
