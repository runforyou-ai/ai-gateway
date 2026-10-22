// Package provider 定义平台预设供应商及写死的协议透传/转换矩阵。
package provider

import "github.com/runforyou-ai/lite-ai-gateway/internal/domain"

// ID 是预设供应商标识。
type ID string

const (
	OpenAI    ID = "openai"
	Anthropic ID = "anthropic"
	Gemini    ID = "gemini"
)

// Spec 描述一个预设供应商的原生协议与可转换对外协议。
type Spec struct {
	ID               ID
	Name             string
	NativeProtocols  []domain.Protocol
	ConvertToNative  map[domain.Protocol]domain.Protocol // 客户端协议 -> 上游原生协议
}

// All 返回全部预设；不支持自定义供应商。
func All() []Spec {
	return []Spec{
		{
			ID:   OpenAI,
			Name: "OpenAI",
			NativeProtocols: []domain.Protocol{
				domain.ProtocolOpenAIChat,
				domain.ProtocolOpenAIResponses,
			},
			ConvertToNative: map[domain.Protocol]domain.Protocol{
				domain.ProtocolAnthropic: domain.ProtocolOpenAIChat,
				domain.ProtocolGemini:    domain.ProtocolOpenAIChat,
			},
		},
		{
			ID:   Anthropic,
			Name: "Anthropic",
			NativeProtocols: []domain.Protocol{
				domain.ProtocolAnthropic,
			},
			ConvertToNative: map[domain.Protocol]domain.Protocol{
				domain.ProtocolOpenAIChat:      domain.ProtocolAnthropic,
				domain.ProtocolOpenAIResponses: domain.ProtocolAnthropic,
				domain.ProtocolGemini:          domain.ProtocolAnthropic,
			},
		},
		{
			ID:   Gemini,
			Name: "Google Gemini",
			NativeProtocols: []domain.Protocol{
				domain.ProtocolGemini,
			},
			ConvertToNative: map[domain.Protocol]domain.Protocol{
				domain.ProtocolOpenAIChat:      domain.ProtocolGemini,
				domain.ProtocolOpenAIResponses: domain.ProtocolGemini,
				domain.ProtocolAnthropic:       domain.ProtocolGemini,
			},
		},
	}
}

// ByID 按标识查找预设；未找到返回 false。
func ByID(id ID) (Spec, bool) {
	for _, spec := range All() {
		if spec.ID == id {
			return spec, true
		}
	}
	return Spec{}, false
}

// Accepts 报告该预设是否接受客户端协议（原生或可转换）。
func (s Spec) Accepts(client domain.Protocol) bool {
	for _, p := range s.NativeProtocols {
		if p == client {
			return true
		}
	}
	_, ok := s.ConvertToNative[client]
	return ok
}

// NeedsConvert 报告客户端协议是否需要转换到上游。
func (s Spec) NeedsConvert(client domain.Protocol) bool {
	for _, p := range s.NativeProtocols {
		if p == client {
			return false
		}
	}
	_, ok := s.ConvertToNative[client]
	return ok
}
