package http

import (
	"encoding/json"
	"errors"
	"github.com/teamcubation/go-items-challenge/internal/adapters/http/presenter"

	customerror "github.com/teamcubation/go-items-challenge/internal/domain/error"
	"github.com/teamcubation/go-items-challenge/internal/domain/user"
	"github.com/teamcubation/go-items-challenge/internal/ports/in"
	"github.com/teamcubation/go-items-challenge/internal/utils"
	"github.com/teamcubation/go-items-challenge/pkg/log"
	"net/http"
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
		w.WriteHeader(http.StatusBadRequest)

		response := presenter.NewErrorResponse(http.StatusBadRequest, "Invalid request payload", map[string]interface{}{
			"error":   err.Error(),
			"context": "Decoding request body",
		})

		json.NewEncoder(w).Encode(response)
		return
	}

	if u.Username == "" || u.Password == "" {
		w.WriteHeader(http.StatusBadRequest)

		response := presenter.NewErrorResponse(http.StatusBadRequest, "Username and password are required", map[string]interface{}{
			"field":   "username",
			"context": "Username and password are required",
		})

		json.NewEncoder(w).Encode(response)
		return
	}

	_, err := h.srv.RegisterUser(ctx, &u)
	if err != nil {
		if errors.Is(err, customerror.ErrUsernameExists) {
			w.WriteHeader(http.StatusBadRequest)

			response := presenter.NewErrorResponse(http.StatusBadRequest, "Username already exists", map[string]interface{}{
				"field":   "username",
				"context": "Username already exists",
			})

			json.NewEncoder(w).Encode(response)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)

		response := presenter.NewErrorResponse(http.StatusInternalServerError, "Internal server error", map[string]interface{}{
			"error":   err.Error(),
			"context": "Creating user",
		})

		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		response := presenter.NewErrorResponse(http.StatusInternalServerError, "Internal server error", map[string]interface{}{
			"error":   err.Error(),
			"context": "Encoding response",
		})

		json.NewEncoder(w).Encode(response)
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

		response := presenter.NewErrorResponse(http.StatusBadRequest, "Invalid request payload", map[string]interface{}{
			"error":   err.Error(),
			"context": "Decoding request body",
		})

		json.NewEncoder(w).Encode(response)
		return
	}

	if creds.Username == "" || creds.Password == "" {
		w.WriteHeader(http.StatusBadRequest)

		response := presenter.NewErrorResponse(http.StatusBadRequest, "Username and password are required", map[string]interface{}{
			"context": "Username and password are required",
		})

		json.NewEncoder(w).Encode(response)
		return
	}
	if err := utils.ValidateStruct(&creds); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		response := presenter.NewErrorResponse(http.StatusBadRequest, "Invalid request payload", map[string]interface{}{
			"error":   err.Error(),
			"context": "missing or invalid fields in the body",
		})

		json.NewEncoder(w).Encode(response)
		return
	}
	token, err := h.srv.Login(ctx, creds)
	if err != nil {
		handleLoginError(w, err, creds.Username)
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]string{"token": token}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		response := presenter.NewErrorResponse(http.StatusInternalServerError, customerror.ErrInternalServer.Error(), map[string]interface{}{
			"error":   err.Error(),
			"context": customerror.ErrEncodingResponse,
		})
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
	}
}

func handleLoginError(w http.ResponseWriter, err error, username string) {
	if errors.Is(err, customerror.ErrUsernameNotFound) {
		response := presenter.NewErrorResponse(http.StatusUnauthorized, customerror.ErrUsernameNotFound.Error(), map[string]interface{}{
			"field":   "username",
			"context": customerror.ErrUsernameNotFound,
		})
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := presenter.NewErrorResponse(http.StatusInternalServerError, customerror.ErrInternalServer.Error(), map[string]interface{}{
		"error":   err.Error(),
		"context": customerror.ErrTokenGeneration,
	})
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(response)
}
