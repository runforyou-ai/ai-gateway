package direct

import (
	"context"
	"errors"

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

// Authorizer 根据请求解析 Cookie 会话；骨架阶段可为 nil（仅 public 可用）。
type Authorizer interface {
	Authenticate(ctx context.Context, r CookieRequest) (auth.Principal, error)
}

// CookieRequest 是鉴权所需的最小请求信息（避免 direct 依赖 net/http 细节时可再收窄）。
type CookieRequest struct {
	CookieHeader string
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
		principal, err := authorizer.Authenticate(ctx, CookieRequest{CookieHeader: meta.Cookie})
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
