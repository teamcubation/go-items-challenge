package http

import (
	_ "context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/teamcubation/go-items-challenge/internal/application"
	"github.com/teamcubation/go-items-challenge/internal/domain/user"
	"github.com/teamcubation/go-items-challenge/internal/ports/in"
	"github.com/teamcubation/go-items-challenge/pkg/log"
)

type AuthHandler struct {
	srv in.AuthService
}

func NewAuthHandler(srv in.AuthService) *AuthHandler {
	return &AuthHandler{srv: srv}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var u user.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if u.Password == "" {
		http.Error(w, "password is required", http.StatusBadRequest)
		return
	}

	if u.Username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	_, err := h.srv.RegisterUser(ctx, &u)
	if err != nil {
		if errors.Is(err, application.ErrUsernameExists) {
			http.Error(w, "username already exists", http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := log.GetFromContext(ctx)
	logger.Info("Entering AuthHandler: Login()")

	var creds user.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	token, err := h.srv.Login(ctx, creds)
	if err != nil {
		if errors.Is(err, application.ErrUsernameNotFound) {
			http.Error(w, "Username not found", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
