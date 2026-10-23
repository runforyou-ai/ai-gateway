package direct

import (
	"context"
	"errors"
	"net/http"

	actionauth "github.com/runforyou-ai/lite-ai-gateway/internal/actions/auth"
	"github.com/runforyou-ai/lite-ai-gateway/internal/actions/installation"
	"github.com/runforyou-ai/lite-ai-gateway/internal/contract"
	"github.com/uptrace/bun"
)

// Ops 是手写业务实现，供生成的 Backend 分发调用。
type Ops struct {
	DB *bun.DB
}

// Health 返回健康状态。
func (o *Ops) Health(context.Context, contract.RequestMeta) (contract.HealthStatus, error) {
	return contract.HealthStatus{OK: true}, nil
}

// InstallationStatus 返回安装状态。
func (o *Ops) InstallationStatus(ctx context.Context, _ contract.RequestMeta) (contract.InstallationStatus, error) {
	installed, siteName, err := installation.Status(ctx, o.DB)
	if err != nil {
		return contract.InstallationStatus{}, err
	}
	return contract.InstallationStatus{Installed: installed, SiteName: siteName}, nil
}

// Register 注册账号。
func (o *Ops) Register(ctx context.Context, _ contract.RequestMeta, input contract.RegisterInput) (contract.RegisterResult, error) {
	out, err := actionauth.Register(ctx, o.DB, actionauth.RegisterInput{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, actionauth.ErrInvalidName):
			return contract.RegisterResult{}, &APIError{
				Status: http.StatusUnprocessableEntity, Message: "invalid name",
				Fields: map[string]string{"name": "invalid"},
			}
		case errors.Is(err, actionauth.ErrInvalidEmail):
			return contract.RegisterResult{}, &APIError{
				Status: http.StatusUnprocessableEntity, Message: "invalid email",
				Fields: map[string]string{"email": "invalid"},
			}
		case errors.Is(err, actionauth.ErrWeakPassword):
			return contract.RegisterResult{}, &APIError{
				Status: http.StatusUnprocessableEntity, Message: "password too weak",
				Fields: map[string]string{"password": "invalid"},
			}
		case errors.Is(err, actionauth.ErrEmailTaken):
			return contract.RegisterResult{}, &APIError{
				Status: http.StatusConflict, Message: "email already registered",
				Fields: map[string]string{"email": "taken"},
			}
		case errors.Is(err, actionauth.ErrRegistrationClosed):
			return contract.RegisterResult{}, &APIError{
				Status: http.StatusForbidden, Message: "registration is closed",
			}
		default:
			return contract.RegisterResult{}, err
		}
	}
	return contract.RegisterResult{Role: string(out.Role)}, nil
}

// NewBackend 组装生成分发层与手写实现。
func NewBackend(db *bun.DB) *Backend {
	return &Backend{Impl: &Ops{DB: db}}
}
