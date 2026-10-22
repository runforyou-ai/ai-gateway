// Package domain 存放领域枚举与常量。
package domain

// Role 是用户角色；仅预设两档，不支持自定义。
type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

// Valid 报告是否为已知角色。
func (r Role) Valid() bool {
	switch r {
	case RoleAdmin, RoleUser:
		return true
	default:
		return false
	}
}
