package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/runforyou-ai/lite-ai-gateway/internal/model"
	"github.com/uptrace/bun"
)

// Me 按用户 id 返回资料。
func Me(ctx context.Context, db bun.IDB, userID int64) (UserInfo, error) {
	var user model.User
	err := db.NewSelect().Model(&user).Where("id = ?", userID).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return UserInfo{}, errors.New("user not found")
	}
	if err != nil {
		return UserInfo{}, err
	}
	return UserInfo{ID: user.ID, Email: user.Email, Role: user.Role}, nil
}
