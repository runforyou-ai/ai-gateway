package direct

import (
	"context"
	"fmt"

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

// Echo 回显输入字符串。
func (o *Ops) Echo(_ context.Context, _ contract.RequestMeta, input contract.EchoInput) (contract.EchoResult, error) {
	return contract.EchoResult{
		Message: fmt.Sprintf("Received, you entered %s", input.Text),
	}, nil
}

// NewBackend 组装生成分发层与手写实现。
func NewBackend(db *bun.DB) *Backend {
	return &Backend{Impl: &Ops{DB: db}}
}
