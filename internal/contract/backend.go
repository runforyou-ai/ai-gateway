package contract

import "context"

//go:generate go run ../tools/contractgen

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

	// InstallationStatus 返回是否已有账号及站点名。
	//contract:route GET /installation/status auth=public
	InstallationStatus(context.Context, RequestMeta) (InstallationStatus, error)

	// Register 注册账号；第一个用户为 admin，其后为 user（需开放注册）。
	//contract:route POST /register status=201 auth=public
	Register(context.Context, RequestMeta, RegisterInput) (RegisterResult, error)

	// ListUsers 分页列出用户（查询参数 _start、_end）。
	//contract:route GET /users auth=admin
	ListUsers(context.Context, RequestMeta) (UserList, error)

	// GetUser 获取单个用户。
	//contract:route GET /users/{id} auth=admin
	GetUser(context.Context, RequestMeta, UserIDInput) (User, error)

	// CreateUser 创建用户。
	//contract:route POST /users status=201 auth=admin
	CreateUser(context.Context, RequestMeta, UserCreateInput) (User, error)

	// UpdateUser 更新用户；password 为空不改密。
	//contract:route PATCH /users/{id} auth=admin
	UpdateUser(context.Context, RequestMeta, UserUpdateInput) (User, error)

	// DeleteUser 删除用户。
	//contract:route DELETE /users/{id} status=204 auth=admin
	DeleteUser(context.Context, RequestMeta, UserIDInput) error

	// InviteUser 邀请用户注册（发信）。
	//contract:route POST /users/invite status=204 auth=admin
	InviteUser(context.Context, RequestMeta, InviteUserInput) error

	// LookupInvite 按令牌查询未使用邀请。
	//contract:route GET /invites/{token} auth=public
	LookupInvite(context.Context, RequestMeta, InviteTokenInput) (InviteLookup, error)

	// AcceptInvite 接受邀请并创建账号。
	//contract:route POST /invites/accept status=201 auth=public
	AcceptInvite(context.Context, RequestMeta, AcceptInviteInput) (User, error)

	// GetSiteSettings 读取站点设置。
	//contract:route GET /settings/site auth=admin
	GetSiteSettings(context.Context, RequestMeta) (SiteSettings, error)

	// UpdateSiteSettings 更新站点设置。
	//contract:route PUT /settings/site auth=admin
	UpdateSiteSettings(context.Context, RequestMeta, SiteSettings) (SiteSettings, error)

	// ListVirtualKeys 分页列出当前用户的虚拟密钥（查询参数 _start、_end）。
	//contract:route GET /virtual-keys auth=user
	ListVirtualKeys(context.Context, RequestMeta) (VirtualKeyList, error)

	// GetVirtualKey 获取当前用户的单个虚拟密钥。
	//contract:route GET /virtual-keys/{id} auth=user
	GetVirtualKey(context.Context, RequestMeta, VirtualKeyIDInput) (VirtualKey, error)

	// CreateVirtualKey 创建虚拟密钥；响应含一次性明文 secret。
	//contract:route POST /virtual-keys status=201 auth=user
	CreateVirtualKey(context.Context, RequestMeta, VirtualKeyCreateInput) (VirtualKeyCreated, error)

	// UpdateVirtualKey 更新名称、禁用状态、预算、重置周期或过期时长。
	//contract:route PATCH /virtual-keys/{id} auth=user
	UpdateVirtualKey(context.Context, RequestMeta, VirtualKeyUpdateInput) (VirtualKey, error)

	// DeleteVirtualKey 删除虚拟密钥。
	//contract:route DELETE /virtual-keys/{id} status=204 auth=user
	DeleteVirtualKey(context.Context, RequestMeta, VirtualKeyIDInput) error
}
