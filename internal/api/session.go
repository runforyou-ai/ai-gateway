package api

import (
	"encoding/json"
	"errors"
	"net/http"

	actionauth "github.com/runforyou-ai/ai-gateway/internal/actions/auth"
	"github.com/runforyou-ai/ai-gateway/internal/auth"
	"github.com/runforyou-ai/ai-gateway/internal/direct"
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
	mux.HandleFunc("PATCH /api/me", h.updateMe)
}

type loginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

type updateMeBody struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

type userBody struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
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
		Email: in.Email, Password: in.Password, Remember: in.Remember,
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
	http.SetCookie(w, auth.SessionCookie(out.Token, h.secure, out.Remember))
	writeJSON(w, http.StatusOK, userBody{
		ID: out.User.ID, Name: out.User.Name, Email: out.User.Email, Role: string(out.User.Role),
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
		ID: info.ID, Name: info.Name, Email: info.Email, Role: string(info.Role),
	}, nil)
}

func (h *sessionHandlers) updateMe(w http.ResponseWriter, r *http.Request) {
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

	var in updateMeBody
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, nil, errBadRequest)
		return
	}

	info, err := actionauth.UpdateProfile(r.Context(), h.db, principal.UserID, actionauth.UpdateProfileInput{
		Name:            in.Name,
		Email:           in.Email,
		CurrentPassword: in.CurrentPassword,
		NewPassword:     in.NewPassword,
	})
	if err != nil {
		switch {
		case errors.Is(err, actionauth.ErrInvalidName):
			writeJSON(w, 0, nil, &direct.APIError{
				Status: http.StatusUnprocessableEntity, Message: "invalid name",
				Fields: map[string]string{"name": "invalid"},
			})
		case errors.Is(err, actionauth.ErrInvalidEmail):
			writeJSON(w, 0, nil, &direct.APIError{
				Status: http.StatusUnprocessableEntity, Message: "invalid email",
				Fields: map[string]string{"email": "invalid"},
			})
		case errors.Is(err, actionauth.ErrEmailTaken):
			writeJSON(w, 0, nil, &direct.APIError{
				Status: http.StatusConflict, Message: "email already registered",
				Fields: map[string]string{"email": "taken"},
			})
		case errors.Is(err, actionauth.ErrWeakPassword):
			writeJSON(w, 0, nil, &direct.APIError{
				Status: http.StatusUnprocessableEntity, Message: "password too weak",
				Fields: map[string]string{"newPassword": "invalid"},
			})
		case errors.Is(err, actionauth.ErrWrongPassword):
			writeJSON(w, 0, nil, &direct.APIError{
				Status: http.StatusUnprocessableEntity, Message: "current password is incorrect",
				Fields: map[string]string{"currentPassword": "invalid"},
			})
		case errors.Is(err, actionauth.ErrUserNotFound):
			writeJSON(w, 0, nil, direct.ErrUnauthorized)
		default:
			writeJSON(w, 0, nil, err)
		}
		return
	}

	writeJSON(w, http.StatusOK, userBody{
		ID: info.ID, Name: info.Name, Email: info.Email, Role: string(info.Role),
	}, nil)
}
