package direct

// APIError 是 /api 可映射为 HTTP 状态与 message/fields 的错误。
type APIError struct {
	Status  int
	Message string
	Fields  map[string]string
}

func (e *APIError) Error() string {
	if e == nil {
		return ""
	}
	return e.Message
}
