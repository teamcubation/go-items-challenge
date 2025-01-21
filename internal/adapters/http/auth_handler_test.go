package http_test

import (
	"bytes"
	"encoding/json"
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

	handler.Register(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid fields username and/or password")

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

	handler.Register(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid fields username and/or password")

	mockService.AssertNotCalled(t, "RegisterUser", mock.Anything, mock.AnythingOfType("*user.User"))
}

func TestAuthHandler_Register_UsernameExists(t *testing.T) {
	mockService := new(mocks.AuthService)
	handler := http2.NewAuthHandler(mockService)

	// Existing username scenario
	inputUser := &user.User{
		Username: "testuser",
		Password: "password123",
	}
	mockService.On("RegisterUser", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil, application.ErrUsernameExists)

	reqBody, _ := json.Marshal(inputUser)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	handler.Register(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "username already exists")

	mockService.AssertCalled(t, "RegisterUser", mock.Anything, mock.AnythingOfType("*user.User"))
}
