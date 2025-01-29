package http_test

import (
	"bytes"
	"encoding/json"
	"github.com/teamcubation/go-items-challenge/internal/adapters/http/middleware"
	"github.com/teamcubation/go-items-challenge/internal/adapters/http/presenter"
	"net/http"
	"net/http/httptest"
	"testing"

	http2 "github.com/teamcubation/go-items-challenge/internal/adapters/http"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/teamcubation/go-items-challenge/internal/domain/user"
	"github.com/teamcubation/go-items-challenge/internal/ports/in/mocks"
)

func executeRequest(req *http.Request, handlerFunc http.HandlerFunc) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler := middleware.ErrorHandlingMiddleware(handlerFunc)
	handler.ServeHTTP(rr, req)
	return rr
}

func TestAuthHandler_Register_Success(t *testing.T) {
	mockService := new(mocks.AuthService)
	handler := http2.NewAuthHandler(mockService)

	// Simule uma requisição válida
	inputUser := &user.User{
		Username: "testuser",
		Password: "password123",
	}
	mockService.On("RegisterUser", mock.Anything, mock.AnythingOfType("*user.User")).Return(inputUser, nil)

	reqBody, _ := json.Marshal(inputUser)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Use executeRequest
	rec := executeRequest(req, handler.Register)

	assert.Equal(t, http.StatusOK, rec.Code)
	var response map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "User created successfully", response["message"])

	mockService.AssertCalled(t, "RegisterUser", mock.Anything, mock.AnythingOfType("*user.User"))
}

func TestAuthHandler_Register_MissingUsername(t *testing.T) {
	mockService := new(mocks.AuthService)
	handler := http2.NewAuthHandler(mockService)

	// Entrada do usuário sem username
	inputUser := &user.User{
		Password: "password123",
	}

	reqBody, _ := json.Marshal(inputUser)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Use executeRequest
	rec := executeRequest(req, handler.Register)

	// Verifique o código de status e a resposta
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var customErr presenter.CustomError
	err := json.Unmarshal(rec.Body.Bytes(), &customErr)
	assert.NoError(t, err)
	assert.Equal(t, "ERR_INVALID_REQUEST_BODY", customErr.Code)
	assert.Equal(t, "Username and password are required", customErr.Message)
	assert.Contains(t, customErr.Details["fields"], "username")
}

func TestAuthHandler_Register_MissingPassword(t *testing.T) {
	mockService := new(mocks.AuthService)
	handler := http2.NewAuthHandler(mockService)

	inputUser := &user.User{
		Username: "testuser",
	}

	reqBody, _ := json.Marshal(inputUser)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rec := executeRequest(req, handler.Register)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var customErr presenter.CustomError
	err := json.Unmarshal(rec.Body.Bytes(), &customErr)
	assert.NoError(t, err)
	assert.Equal(t, "ERR_INVALID_REQUEST_BODY", customErr.Code)
	assert.Equal(t, "Username and password are required", customErr.Message)
	assert.Contains(t, customErr.Details["fields"], "password")
}

func TestAuthHandler_Register_UsernameExists(t *testing.T) {
	mockService := new(mocks.AuthService)
	handler := http2.NewAuthHandler(mockService)

	inputUser := &user.User{
		Username: "testuser",
		Password: "password123",
	}

	mockService.On("RegisterUser", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil, presenter.ErrUsernameExists)

	reqBody, _ := json.Marshal(inputUser)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Use executeRequest
	rec := executeRequest(req, handler.Register)

	assert.Equal(t, http.StatusConflict, rec.Code)

	var customErr presenter.CustomError
	err := json.Unmarshal(rec.Body.Bytes(), &customErr)
	assert.NoError(t, err)
	assert.Equal(t, "ERR_USERNAME_EXISTS", customErr.Code)
	assert.Equal(t, "Username already exists", customErr.Message)
}

func TestAuthHandler_Login_Success(t *testing.T) {
	mockService := new(mocks.AuthService)
	handler := http2.NewAuthHandler(mockService)

	creds := user.Credentials{
		Username: "validuser",
		Password: "password123",
	}

	mockService.On("Login", mock.Anything, creds).Return("valid_token", nil)

	reqBody, _ := json.Marshal(creds)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.Login(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "valid_token", response["token"])
}
