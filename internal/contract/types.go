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

// VirtualKey 是管理台虚拟密钥视图（不含明文）。
type VirtualKey struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Prefix      string `json:"prefix"`
	Disabled    bool   `json:"disabled"`
	MaxBudget   string `json:"maxBudget"`
	Spent       string `json:"spent"`
	BudgetReset string `json:"budgetReset"`
	ExpiresIn   string `json:"expiresIn"`
	ExpiresAt   string `json:"expiresAt"`
	LastUsedAt  string `json:"lastUsedAt"`
	CreatedAt   string `json:"createdAt"`
}

// VirtualKeyList 是虚拟密钥分页列表。
type VirtualKeyList struct {
	Items []VirtualKey `json:"items"`
	Total int64        `json:"total"`
}

// VirtualKeyIDInput 按 id 定位虚拟密钥。
type VirtualKeyIDInput struct {
	ID int64 `json:"id"`
}

// VirtualKeyCreateInput 是创建虚拟密钥参数；限制类字段空表示不限。
type VirtualKeyCreateInput struct {
	Name        string `json:"name"`
	MaxBudget   string `json:"maxBudget"`
	ExpiresIn   string `json:"expiresIn"`
	BudgetReset string `json:"budgetReset"`
}

// VirtualKeyCreated 是创建响应，含一次性明文。
type VirtualKeyCreated struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Prefix      string `json:"prefix"`
	Disabled    bool   `json:"disabled"`
	MaxBudget   string `json:"maxBudget"`
	Spent       string `json:"spent"`
	BudgetReset string `json:"budgetReset"`
	ExpiresIn   string `json:"expiresIn"`
	ExpiresAt   string `json:"expiresAt"`
	LastUsedAt  string `json:"lastUsedAt"`
	CreatedAt   string `json:"createdAt"`
	Secret      string `json:"secret"`
}

// VirtualKeyUpdateInput 是更新虚拟密钥参数；限制类字段空表示清除限制。
type VirtualKeyUpdateInput struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Disabled    bool   `json:"disabled"`
	MaxBudget   string `json:"maxBudget"`
	ExpiresIn   string `json:"expiresIn"`
	BudgetReset string `json:"budgetReset"`
}

// CatalogModel 是管理台模型目录视图；对外协议由 type 在代码中推导，不单独落库或下发。
type CatalogModel struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Type        string `json:"type"`
	Disabled    bool   `json:"disabled"`
	CreatedAt   string `json:"createdAt"`
}

// CatalogModelList 是模型目录分页列表。
type CatalogModelList struct {
	Items []CatalogModel `json:"items"`
	Total int64          `json:"total"`
}

// CatalogModelIDInput 按 id 定位目录模型。
type CatalogModelIDInput struct {
	ID int64 `json:"id"`
}

// CatalogModelCreateInput 是创建目录模型参数。
type CatalogModelCreateInput struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Type        string `json:"type"`
	Disabled    bool   `json:"disabled"`
}

// CatalogModelUpdateInput 是更新目录模型参数。
type CatalogModelUpdateInput struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Type        string `json:"type"`
	Disabled    bool   `json:"disabled"`
}
