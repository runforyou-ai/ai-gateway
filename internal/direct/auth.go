package direct

import (
	"context"
	"errors"
	"strings"

	"github.com/runforyou-ai/lite-ai-gateway/internal/auth"
	"github.com/runforyou-ai/lite-ai-gateway/internal/contract"
	"github.com/runforyou-ai/lite-ai-gateway/internal/domain"
)

var (
	// ErrUnauthorized 表示未认证。
	ErrUnauthorized = errors.New("unauthorized")
	// ErrForbidden 表示已认证但无权限。
	ErrForbidden = errors.New("forbidden")
)

// Authorizer 校验 JWT；骨架阶段可为 nil（仅 public 可用）。
type Authorizer interface {
	Authenticate(ctx context.Context, authorization string) (auth.Principal, error)
}

// Options 配置分发层依赖。
type Options struct {
	Authorizer Authorizer
}

var authorizer Authorizer

// SetAuthorizer 设置全局鉴权器（serve 启动时注入）。
func SetAuthorizer(a Authorizer) {
	authorizer = a
}

func authorize(ctx context.Context, meta contract.RequestMeta, mode string) error {
	switch mode {
	case "public":
		return nil
	case "user", "admin":
		if authorizer == nil {
			return ErrUnauthorized
		}
		principal, err := authorizer.Authenticate(ctx, meta.Authorization)
		if err != nil {
			return ErrUnauthorized
		}
		if mode == "admin" && principal.Role != domain.RoleAdmin {
			return ErrForbidden
		}
		return nil
	default:
		return ErrUnauthorized
	}
}

// BearerToken 从 Authorization 头取出 bearer token。
func BearerToken(authorization string) string {
	const prefix = "Bearer "
	if strings.HasPrefix(authorization, prefix) {
		return strings.TrimSpace(authorization[len(prefix):])
	}
	return ""
}
