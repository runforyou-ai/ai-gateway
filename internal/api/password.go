package api

import (
	"encoding/json"
	"errors"
	"net/http"

	actionauth "github.com/runforyou-ai/ai-gateway/internal/actions/auth"
	"github.com/runforyou-ai/ai-gateway/internal/direct"
	"github.com/runforyou-ai/ai-gateway/internal/email"
	"github.com/uptrace/bun"
)

type passwordHandlers struct {
	db        *bun.DB
	mailer    email.Sender
	publicURL string
}

func registerPasswordRoutes(mux *http.ServeMux, db *bun.DB, mailer email.Sender, publicURL string) {
	h := &passwordHandlers{db: db, mailer: mailer, publicURL: publicURL}
	mux.HandleFunc("POST /api/forgot-password", h.forgotPassword)
	mux.HandleFunc("POST /api/reset-password", h.resetPassword)
}

type forgotPasswordBody struct {
	Email string `json:"email"`
}

type resetPasswordBody struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

func (h *passwordHandlers) forgotPassword(w http.ResponseWriter, r *http.Request) {
	var in forgotPasswordBody
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, nil, errBadRequest)
		return
	}
	err := actionauth.RequestPasswordReset(r.Context(), h.db, h.mailer, h.publicURL, in.Email)
	if err != nil {
		switch {
		case errors.Is(err, actionauth.ErrInvalidEmail):
			writeJSON(w, 0, nil, &direct.APIError{
				Status: http.StatusUnprocessableEntity, Message: "invalid email",
				Fields: map[string]string{"email": "invalid"},
			})
		case errors.Is(err, actionauth.ErrEmailNotConfigured):
			writeJSON(w, 0, nil, &direct.APIError{
				Status: http.StatusServiceUnavailable, Message: "email not configured",
			})
		default:
			writeJSON(w, 0, nil, err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *passwordHandlers) resetPassword(w http.ResponseWriter, r *http.Request) {
	var in resetPasswordBody
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, nil, errBadRequest)
		return
	}
	err := actionauth.ResetPassword(r.Context(), h.db, in.Token, in.Password)
	if err != nil {
		switch {
		case errors.Is(err, actionauth.ErrInvalidResetToken):
			writeJSON(w, 0, nil, &direct.APIError{
				Status: http.StatusUnprocessableEntity, Message: "invalid or expired reset token",
				Fields: map[string]string{"token": "invalid"},
			})
		case errors.Is(err, actionauth.ErrWeakPassword):
			writeJSON(w, 0, nil, &direct.APIError{
				Status: http.StatusUnprocessableEntity, Message: "password too weak",
				Fields: map[string]string{"password": "invalid"},
			})
		case errors.Is(err, actionauth.ErrUserNotFound):
			writeJSON(w, 0, nil, &direct.APIError{
				Status: http.StatusUnprocessableEntity, Message: "invalid or expired reset token",
				Fields: map[string]string{"token": "invalid"},
			})
		default:
			writeJSON(w, 0, nil, err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
