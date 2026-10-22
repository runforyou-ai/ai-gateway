package domain

// Protocol 是对外 / 上游 API 协议标识。
type Protocol string

const (
	ProtocolOpenAIChat      Protocol = "openai.chat"
	ProtocolOpenAIResponses Protocol = "openai.responses"
	ProtocolAnthropic       Protocol = "anthropic.messages"
	ProtocolGemini          Protocol = "gemini.generate_content"
)
