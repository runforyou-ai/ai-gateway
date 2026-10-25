package contract

// HealthStatus 是健康检查结果。
type HealthStatus struct {
	OK bool `json:"ok"`
}

// InstallationStatus 描述是否已有账号（首个注册用户即为管理员）。
type InstallationStatus struct {
	Installed        bool   `json:"installed"`
	SiteName         string `json:"siteName"`
	RegistrationOpen bool   `json:"registrationOpen"`
}

// RegisterInput 是注册请求。
type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterResult 是注册响应。
type RegisterResult struct {
	Role string `json:"role"`
}

// User 是管理台用户视图（不含密码）。
type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Disabled bool   `json:"disabled"`
}

// UserList 是用户分页列表。
type UserList struct {
	Items []User `json:"items"`
	Total int64  `json:"total"`
}

// UserIDInput 按 id 定位用户。
type UserIDInput struct {
	ID int64 `json:"id"`
}

// UserCreateInput 是管理员创建用户参数。
type UserCreateInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Disabled bool   `json:"disabled"`
}

// UserUpdateInput 是管理员更新用户参数；password 为空表示不改密码。
type UserUpdateInput struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Disabled bool   `json:"disabled"`
}

// InviteUserInput 是管理员邀请用户参数。
type InviteUserInput struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

// InviteTokenInput 按令牌定位邀请。
type InviteTokenInput struct {
	Token string `json:"token"`
}

// InviteLookup 是公开查询邀请的结果。
type InviteLookup struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

// AcceptInviteInput 是接受邀请参数。
type AcceptInviteInput struct {
	Token    string `json:"token"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// SiteSettings 是站点设置。
type SiteSettings struct {
	Name             string `json:"name"`
	RegistrationOpen bool   `json:"registrationOpen"`
}
