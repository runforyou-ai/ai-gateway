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

// defaultFromName 在配置仅给邮箱地址时补上的显示名，避免收件箱只显示本地部分（如 hi）。
const defaultFromName = "AI Gateway"

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
	return newResend(strings.TrimSpace(opts.ResendAPIKey), normalizeFrom(opts.From))
}

// normalizeFrom 保证 From 带有显示名。裸地址会补 defaultFromName；已写「名称 <邮箱>」则保留名称。
func normalizeFrom(from string) string {
	from = strings.TrimSpace(from)
	if from == "" {
		return ""
	}
	addr, err := mail.ParseAddress(from)
	if err != nil {
		return from
	}
	if strings.TrimSpace(addr.Name) == "" {
		addr.Name = defaultFromName
	}
	return addr.String()
}

type notConfigured struct{}

func (notConfigured) Send(context.Context, Message) error {
	return ErrNotConfigured
}
