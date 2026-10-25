package contract

import "context"

// RequestMeta 是契约方法共用的请求元数据。
type RequestMeta struct {
	// Cookie 是原始 Cookie 头，管理台会话鉴权使用。
	Cookie string
	// RemoteAddr 是客户端地址。
	RemoteAddr string
	// Query 是 URL 原始查询串（不含 ?）。
	Query string
	// ActorUserID 由鉴权填充；未认证时为 0。
	ActorUserID int64
	// ActorRole 由鉴权填充；未认证时为空。
	ActorRole string
}

type metaKey struct{}

// WithMeta 将 meta 放入 context。
func WithMeta(ctx context.Context, meta RequestMeta) context.Context {
	return context.WithValue(ctx, metaKey{}, meta)
}

// MetaFrom 从 context 取 meta。
func MetaFrom(ctx context.Context) RequestMeta {
	meta, _ := ctx.Value(metaKey{}).(RequestMeta)
	return meta
}
