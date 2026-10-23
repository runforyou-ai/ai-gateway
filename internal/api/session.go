package api

import (
	"encoding/json"
	"errors"
	"net/http"

	actionauth "github.com/runforyou-ai/lite-ai-gateway/internal/actions/auth"
	"github.com/runforyou-ai/lite-ai-gateway/internal/auth"
	"github.com/runforyou-ai/lite-ai-gateway/internal/direct"
	"github.com/uptrace/bun"
)

type sessionHandlers struct {
	db     *bun.DB
	secure bool
}

func registerSessionRoutes(mux *http.ServeMux, db *bun.DB, secure bool) {
	h := &sessionHandlers{db: db, secure: secure}
	mux.HandleFunc("POST /api/login", h.login)
	mux.HandleFunc("POST /api/logout", h.logout)
	mux.HandleFunc("GET /api/me", h.me)
}

type loginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userBody struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (h *sessionHandlers) login(w http.ResponseWriter, r *http.Request) {
	var in loginBody
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, nil, errBadRequest)
		return
	}
	out, err := actionauth.Login(r.Context(), h.db, actionauth.LoginInput{
		Email: in.Email, Password: in.Password,
	})
	if errors.Is(err, actionauth.ErrInvalidCredentials) {
		writeJSON(w, 0, nil, &direct.APIError{
			Status: http.StatusUnauthorized, Message: "invalid email or password",
		})
		return
	}
	if err != nil {
		writeJSON(w, 0, nil, err)
		return
	}
	http.SetCookie(w, auth.SessionCookie(out.Token, h.secure))
	writeJSON(w, http.StatusOK, userBody{
		ID: out.User.ID, Email: out.User.Email, Role: string(out.User.Role),
	}, nil)
}

func (h *sessionHandlers) logout(w http.ResponseWriter, r *http.Request) {
	if err := actionauth.Logout(r.Context(), h.db, r.Header.Get("Cookie")); err != nil {
		writeJSON(w, 0, nil, err)
		return
	}
	http.SetCookie(w, auth.ClearSessionCookie(h.secure))
	w.WriteHeader(http.StatusNoContent)
}

func (h *sessionHandlers) me(w http.ResponseWriter, r *http.Request) {
	authorizer := direct.GetAuthorizer()
	if authorizer == nil {
		writeJSON(w, 0, nil, direct.ErrUnauthorized)
		return
	}
	principal, err := authorizer.Authenticate(r.Context(), direct.CookieRequest{
		CookieHeader: r.Header.Get("Cookie"),
	})
	if err != nil {
		writeJSON(w, 0, nil, direct.ErrUnauthorized)
		return
	}
	info, err := actionauth.Me(r.Context(), h.db, principal.UserID)
	if err != nil {
		writeJSON(w, 0, nil, direct.ErrUnauthorized)
		return
	}
	writeJSON(w, http.StatusOK, userBody{
		ID: info.ID, Email: info.Email, Role: string(info.Role),
	}, nil)
}
