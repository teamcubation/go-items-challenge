package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/teamcubation/go-items-challenge/internal/adapters/http/presenter"
	errs "github.com/teamcubation/go-items-challenge/internal/domain/errors"
	"github.com/teamcubation/go-items-challenge/internal/domain/user"
	"github.com/teamcubation/go-items-challenge/internal/ports/in"
	"github.com/teamcubation/go-items-challenge/internal/utils"
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
	w.Header().Set("Content-Type", "application/json")

	var u user.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		response := presenter.NewApiError(http.StatusBadRequest, "Invalid request payload", map[string]interface{}{
			"error":   err.Error(),
			"context": "Decoding request body",
		})

		json.NewEncoder(w).Encode(response)
		return
	}

	if u.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(presenter.NewApiError(http.StatusBadRequest, "Password is required", map[string]interface{}{
			"filed":   "password",
			"context": "Validating request body",
		}))
		return
	}

	if u.Username == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(presenter.NewApiError(http.StatusBadRequest, "Username is required", map[string]interface{}{
			"filed":   "username",
			"context": "Validating request body",
		}))
		return
	}

	_, err := h.srv.RegisterUser(ctx, &u)
	if err != nil {
		if errors.Is(err, errs.ErrUsernameExists) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(presenter.NewApiError(http.StatusBadRequest, "Username already exists", map[string]interface{}{
				"field":   "username",
				"context": "Error registering user",
			}))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(presenter.NewApiError(http.StatusInternalServerError, "Internal server error", map[string]interface{}{
			"error":   "Internal server error",
			"context": "Registering user",
		}))
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(presenter.NewApiError(http.StatusInternalServerError, "Internal server error", map[string]interface{}{
			"error":   "Internal server error",
			"context": "Encoding response",
		}))
		return
	}
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
		w.WriteHeader(http.StatusBadRequest)

		response := presenter.NewApiError(http.StatusBadRequest, "Invalid request payload", map[string]interface{}{
			"error":   err.Error(),
			"context": "Decoding request body",
		})

		json.NewEncoder(w).Encode(response)
		return
	}

	if creds.Username == "" || creds.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(presenter.NewApiError(http.StatusBadRequest, "Username and password must not be empty", map[string]interface{}{
			"field":   "username/password",
			"context": "Checking input data",
		}))
		return
	}

	if err := utils.ValidateStruct(&creds); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(presenter.NewApiError(http.StatusBadRequest, "Username and password must not be empty", map[string]interface{}{
			"field":   "username/password",
			"context": "Checking input data",
		}))
		return
	}

	token, err := h.srv.Login(ctx, creds)
	if err != nil {
		if errors.Is(err, errs.ErrUsernameNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(presenter.NewApiError(http.StatusBadRequest, "Username already exists", map[string]interface{}{
				"field":   "username",
				"context": "Error authenticating user",
			}))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(presenter.NewApiError(http.StatusInternalServerError, "Internal server error", map[string]interface{}{
			"error":   "Internal server error",
			"context": "Authenticating user",
		}))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"token": token}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(presenter.NewApiError(http.StatusInternalServerError, "Internal server error", map[string]interface{}{
			"error":   "Internal server error",
			"context": "Encoding response",
		}))
	}
}
