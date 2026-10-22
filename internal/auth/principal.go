// Package auth 提供 JWT 与会话相关原语。
package auth

import "github.com/runforyou-ai/lite-ai-gateway/internal/domain"

// Principal 是已认证主体。
type Principal struct {
	UserID int64
	Role   domain.Role
}
