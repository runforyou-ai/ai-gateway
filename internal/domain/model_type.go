package domain

// ModelType 是目录模型类型；同类型内协议可互相转换。
type ModelType string

const (
	ModelTypeChat ModelType = "chat"
)

// AllModelTypes 返回管理台可选的模型类型。
func AllModelTypes() []ModelType {
	return []ModelType{ModelTypeChat}
}

// Valid 报告是否为已知模型类型。
func (t ModelType) Valid() bool {
	switch t {
	case ModelTypeChat:
		return true
	default:
		return false
	}
}

// Protocols 返回该类型下全部预设对外协议（由代码约定，避免建模型时漏选）。
func (t ModelType) Protocols() []Protocol {
	switch t {
	case ModelTypeChat:
		return []Protocol{
			ProtocolOpenAIChat,
			ProtocolOpenAIResponses,
			ProtocolAnthropic,
			ProtocolGemini,
		}
	default:
		return nil
	}
}
