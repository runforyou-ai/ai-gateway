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

type errorBody struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, body any, err error) {
	if err != nil {
		code := http.StatusInternalServerError
		msg := "internal error"
		switch {
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(errorBody{Error: msg})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if body == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(body)
}
