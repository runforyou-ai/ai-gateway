package contract

// HealthStatus 是健康检查结果。
type HealthStatus struct {
	OK bool `json:"ok"`
}

// InstallationStatus 描述是否已有账号（首个注册用户即为管理员）。
type InstallationStatus struct {
	Installed bool   `json:"installed"`
	SiteName  string `json:"siteName"`
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
