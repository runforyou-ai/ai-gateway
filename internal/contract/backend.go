package contract

import "context"

//go:generate go run github.com/runforyou-ai/lite-ai-gateway/internal/tools/contractgen

// Backend 是管理 API 的唯一契约源。
//
// 每个方法必须携带一条 contract:route 指令：
//
//	contract:route <HTTP方法> <路径> [status=201] [auth=public|user|admin]
//
// 路径相对 /api。auth 默认为 user。
type Backend interface {
	// Health 返回服务健康状态。
	//contract:route GET /health auth=public
	Health(context.Context, RequestMeta) (HealthStatus, error)

	// InstallationStatus 返回是否已安装及站点名。
	//contract:route GET /installation/status auth=public
	InstallationStatus(context.Context, RequestMeta) (InstallationStatus, error)
}
