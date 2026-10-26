package domain

// Protocol 是对外 / 上游 API 协议标识（预设，不可自定义）。
type Protocol string

const (
	ProtocolOpenAIChat      Protocol = "openai.chat"
	ProtocolOpenAIResponses Protocol = "openai.responses"
	ProtocolAnthropic       Protocol = "anthropic.messages"
	ProtocolGemini          Protocol = "gemini.generate_content"
)

// AllProtocols 返回全部预设协议。
func AllProtocols() []Protocol {
	out := make([]Protocol, 0, 8)
	for _, t := range AllModelTypes() {
		out = append(out, t.Protocols()...)
	}
	return out
}

// Valid 报告是否为已知协议。
func (p Protocol) Valid() bool {
	_, ok := p.ModelType()
	return ok
}

// ModelType 返回协议所属模型类型；未知协议返回 false。
func (p Protocol) ModelType() (ModelType, bool) {
	for _, t := range AllModelTypes() {
		for _, q := range t.Protocols() {
			if p == q {
				return t, true
			}
		}
	}
	return "", false
}
