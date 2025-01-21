package http_test

import (
	"bytes"
	"encoding/json"
	errs "github.com/teamcubation/go-items-challenge/errors"
	"net/http"
	"net/http/httptest"
	"testing"

	http2 "github.com/teamcubation/go-items-challenge/internal/adapters/http"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/teamcubation/go-items-challenge/internal/application"
	"github.com/teamcubation/go-items-challenge/internal/domain/user"
	"github.com/teamcubation/go-items-challenge/internal/ports/in/mocks"
)

func TestAuthHandler_Register_Success(t *testing.T) {
	mockService := new(mocks.AuthService)
	handler := http2.NewAuthHandler(mockService)

	// Simulate a valid request
	inputUser := &user.User{
		Username: "testuser",
		Password: "password123",
	}
	mockService.On("RegisterUser", mock.Anything, mock.AnythingOfType("*user.User")).Return(inputUser, nil)

	reqBody, _ := json.Marshal(inputUser)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	handler.Register(rec, req)

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

	// User input missing username
	inputUser := &user.User{
		Password: "password123",
	}

	reqBody, _ := json.Marshal(inputUser)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	// Capturar o panic
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(*errs.CustomError); ok && err.StatusCode == http.StatusBadRequest {
				assert.Contains(t, err.Message, "Username is required")
			} else {
				t.Errorf("Unexpected panic: %v", r)
			}
		}
	}()

	handler.Register(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockService.AssertNotCalled(t, "RegisterUser", mock.Anything, mock.AnythingOfType("*user.User"))
}

func TestAuthHandler_Register_MissingPassword(t *testing.T) {
	mockService := new(mocks.AuthService)
	handler := http2.NewAuthHandler(mockService)

	// User input missing password
	inputUser := &user.User{
		Username: "testuser",
	}

	reqBody, _ := json.Marshal(inputUser)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	// Capturar o panic
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(*errs.CustomError); ok && err.StatusCode == http.StatusBadRequest {
				assert.Contains(t, err.Message, "Password is required")
			} else {
				t.Errorf("Unexpected panic: %v", r)
			}
		}
	}()

	handler.Register(rec, req)

	// Opcional: Você pode verificar o código HTTP e a resposta mesmo que não haja pânico
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAuthHandler_Register_UsernameExists(t *testing.T) {
	mockService := new(mocks.AuthService)
	handler := http2.NewAuthHandler(mockService)

	// Simule a existência do nome de usuário
	inputUser := &user.User{
		Username: "testuser",
		Password: "password123",
	}

	mockService.On("RegisterUser", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil, application.ErrUsernameExists)

	reqBody, _ := json.Marshal(inputUser)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	// Capturar o panic
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(*errs.CustomError); ok && err.StatusCode == http.StatusBadRequest {
				assert.Contains(t, err.Message, "Username already exists")
			} else {
				t.Errorf("Unexpected panic: %v", r)
			}
		}
	}()

	handler.Register(rec, req)

	// Verifique que o mockService foi chamado como esperado
	mockService.AssertCalled(t, "RegisterUser", mock.Anything, mock.AnythingOfType("*user.User"))
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

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	mockService := new(mocks.AuthService)
	handler := http2.NewAuthHandler(mockService)

	creds := user.Credentials{
		Username: "invaliduser",
		Password: "wrongpassword",
	}

	mockService.On("Login", mock.Anything, creds).Return("", application.ErrUsernameNotFound)

	reqBody, _ := json.Marshal(creds)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(*errs.CustomError); ok && err.StatusCode == http.StatusUnauthorized {
				assert.Contains(t, err.Message, "Invalid credentials")
			} else {
				t.Errorf("Unexpected panic: %v", r)
			}
		}
	}()

	handler.Login(rec, req)

	mockService.AssertCalled(t, "Login", mock.Anything, creds)
}
