package api

import (
	"encoding/json"
	"imdb/internal/security"
	"imdb/internal/service"
	"log/slog"
	"net/http"
	"time"
)

type UsersHandler struct {
	s *service.UserService
	l *slog.Logger
	p *security.TokenService
}

func NewUsersHandler(s *service.UserService, l *slog.Logger, p *security.TokenService) *UsersHandler {
	return &UsersHandler{s: s, l: l, p: p}
}

func (h *UsersHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var payload LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.l.Error("UsersHandler.SignIn", "decoding error", err, "email", payload.Email)
		ErrorResponse(w, http.StatusBadRequest, ErrBadPayload.Error())
		return
	}

	user, err := h.s.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		h.l.Error("UsersHandler.SignIn", "db error", err, "email", payload.Email)
		ErrorResponse(w, http.StatusBadRequest, ErrBadCredentials.Error())
		return
	}

	if identicalPasswords := security.CheckPassword(payload.Password, user.PasswordHash); !identicalPasswords {
		h.l.Error("UsersHandler.SignIn", "password comparison error", err, "email", payload.Email)
		ErrorResponse(w, http.StatusBadRequest, ErrBadCredentials.Error())
		return
	}

	// we'd not make this much ttl for web token normally
	// for local usage only
	token := h.p.Generate(user.ID.Hex(), 24*time.Hour)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
	if err != nil {
		h.l.Error("UsersHandler.SignIn", "encoding error", err, "email", payload.Email)
		ErrorResponse(w, http.StatusInternalServerError, ErrInternal.Error())
		return
	}
}

func (h *UsersHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	// naming might be confusing
	// but creating a separate struct for same data feels wrong
	var payload LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.l.Error("UsersHandler.SingUp", "decoding error", err, "email", payload.Email)
		ErrorResponse(w, http.StatusBadRequest, ErrBadPayload.Error())
		return
	}

	if err := h.s.Create(r.Context(), payload.Email, payload.Password); err != nil {
		h.l.Error("UsersHandler.SingUp", "creating error", err, "email", payload.Email)
		ErrorResponse(w, http.StatusInternalServerError, ErrInternal.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}
