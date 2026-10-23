package contract

// HealthStatus 是健康检查结果。
type HealthStatus struct {
	OK bool `json:"ok"`
}

// InstallationStatus 描述是否已完成首次安装（至少一名用户）。
type InstallationStatus struct {
	Installed bool   `json:"installed"`
	SiteName  string `json:"siteName"`
}

// EchoInput 是回显请求。
type EchoInput struct {
	Text string `json:"text"`
}

// EchoResult 是回显响应。
type EchoResult struct {
	Message string `json:"message"`
}
