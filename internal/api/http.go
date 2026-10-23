package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/runforyou-ai/lite-ai-gateway/internal/contract"
	"github.com/runforyou-ai/lite-ai-gateway/internal/direct"
)

func requestMeta(r *http.Request) contract.RequestMeta {
	return contract.RequestMeta{
		Cookie:     r.Header.Get("Cookie"),
		RemoteAddr: r.RemoteAddr,
	}
}

func decodeJSONBody(r *http.Request, dst any) error {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("%w: %v", errBadRequest, err)
	}
	return nil
}

var errBadRequest = errors.New("bad request")

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

type errorBody struct {
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func writeAPIError(w http.ResponseWriter, status int, message string, fields map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorBody{Message: message, Fields: fields})
}

func writeJSON(w http.ResponseWriter, status int, body any, err error) {
	if err != nil {
		code := http.StatusInternalServerError
		msg := "internal error"
		var fields map[string]string
		var apiErr *APIError
		switch {
		case errors.As(err, &apiErr):
			code = apiErr.Status
			msg = apiErr.Message
			fields = apiErr.Fields
		case errors.Is(err, direct.ErrUnauthorized):
			code = http.StatusUnauthorized
			msg = "unauthorized"
		case errors.Is(err, direct.ErrForbidden):
			code = http.StatusForbidden
			msg = "forbidden"
		case errors.Is(err, errBadRequest):
			code = http.StatusBadRequest
			msg = "bad request"
		default:
			slog.Error("api request failed", "err", err)
		}
		writeAPIError(w, code, msg, fields)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if body == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(body)
}
