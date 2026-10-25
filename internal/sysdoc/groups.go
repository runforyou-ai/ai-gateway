package sysdoc

// 设置组名。
const (
	SettingsSite = "site"
)

// 状态组名。
const (
	StateInstallation = "installation"
)

// Site 是站点功能设置。
type Site struct {
	Name               string `json:"name"`
	RegistrationOpen   bool   `json:"registrationOpen"`
}

// DefaultSite 返回安装时的站点默认值。
func DefaultSite() Site {
	return Site{
		Name:             "AI Gateway",
		RegistrationOpen: false,
	}
}

// Installation 记录首次安装完成状态。
type Installation struct {
	Installed bool `json:"installed"`
}
