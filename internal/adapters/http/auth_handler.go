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

// Register Registra um novo usuário
// @Summary Registra um novo usuário
// @Description Cria uma nova conta de usuário com os dados fornecidos no corpo da requisição
// @Tags auth
// @Accept json
// @Produce json
// @Param user body user.User true "Informações do usuário"
// @Success 200 {object} map[string]string "Usuário criado com sucesso"
// @Failure 400 {string} string "Username ou senha é obrigatória"
// @Failure 500 {string} string "Erro interno do servidor"
// @Router /register [post]
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

// Login Autentica um usuário
// @Summary Autentica um usuário
// @Description Autentica um usuário com as credenciais fornecidas no corpo da requisição
// @Tags auth
// @Accept json
// @Produce json
// @Param user body user.Credentials true "Credenciais do usuário"
// @Success 200 {object} map[string]string "Token de autenticação"
// @Failure 400 {string} string "Credenciais inválidas"
// @Failure 401 {string} string "Usuário não encontrado"
// @Failure 500 {string} string "Erro interno do servidor"
// @Router /login [post]
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
