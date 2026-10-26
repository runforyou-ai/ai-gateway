// Package email 提供部署配置驱动的发信能力；当前仅支持 Resend。
package email

import (
	"context"
	"errors"
	"net/mail"
	"strings"
)

// ErrNotConfigured 表示未配置发信（缺 from 或 Resend API Key）。
var ErrNotConfigured = errors.New("email not configured")

// Message 是一封待发邮件。
type Message struct {
	To      string
	Subject string
	HTML    string
	Text    string
}

// Sender 发送邮件。
type Sender interface {
	Send(ctx context.Context, msg Message) error
}

// Options 构造发信器所需的部署项。
type Options struct {
	From         string
	FromName     string
	ResendAPIKey string
}

// Configured 报告 Options 是否足以发信。
func (o Options) Configured() bool {
	return strings.TrimSpace(o.From) != "" && strings.TrimSpace(o.ResendAPIKey) != ""
}

// New 按 Options 返回 Sender；未配齐时返回的 Sender 在 Send 时返回 ErrNotConfigured。
func New(opts Options) Sender {
	if !opts.Configured() {
		return notConfigured{}
	}
	return newResend(strings.TrimSpace(opts.ResendAPIKey), formatFrom(opts.From, opts.FromName))
}

// formatFrom 组合地址与显示名。from_name 优先；否则保留 from 里已写的「名称 <邮箱>」。
// 简单显示名输出为 Name <addr>（与 Resend 文档一致）；net/mail 默认会加引号，Resend 常弄丢显示名。
func formatFrom(from, fromName string) string {
	from = strings.TrimSpace(from)
	fromName = strings.TrimSpace(fromName)
	if from == "" {
		return ""
	}
	addr, err := mail.ParseAddress(from)
	if err != nil {
		return from
	}
	name := strings.TrimSpace(addr.Name)
	if fromName != "" {
		name = fromName
	}
	if name == "" {
		return addr.Address
	}
	if displayNameNeedsQuotes(name) {
		return (&mail.Address{Name: name, Address: addr.Address}).String()
	}
	return name + " <" + addr.Address + ">"
}

func displayNameNeedsQuotes(name string) bool {
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
		case r == ' ' || r == '-' || r == '_' || r == '.':
		default:
			return true
		}
	}
	return false
}

type notConfigured struct{}

func (notConfigured) Send(context.Context, Message) error {
	return ErrNotConfigured
}
