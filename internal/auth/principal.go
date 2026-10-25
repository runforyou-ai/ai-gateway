// Package auth 提供 Cookie 会话与主体原语。
package auth

import "github.com/runforyou-ai/ai-gateway/internal/domain"

// Principal 是已认证主体（由有效 session 解析而来）。
type Principal struct {
	UserID    int64
	SessionID int64
	Role      domain.Role
}
