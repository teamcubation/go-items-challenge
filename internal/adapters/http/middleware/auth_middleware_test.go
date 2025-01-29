package middleware_test

// Inserção de partes ausentes e verificação dos imports
import (
	"encoding/json"
	"github.com/teamcubation/go-items-challenge/internal/adapters/http/presenter"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/teamcubation/go-items-challenge/internal/adapters/http/middleware"
)

func createToken(userID int, secret []byte) (string, error) {
	claims := &middleware.Claims{
		UserID:         userID,
		StandardClaims: jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	token, err := createToken(123, middleware.JwtKey)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()

	handler := middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middleware.UserContextKey).(int)
		assert.True(t, ok)
		assert.Equal(t, 123, userID)
		w.WriteHeader(http.StatusOK)
	}))

	// Adicionar o ErrorHandlingMiddleware
	finalHandler := middleware.ErrorHandlingMiddleware(handler)

	finalHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handlerWithAuth := middleware.AuthMiddleware(testHandler)

	// Adicionar o ErrorHandlingMiddleware
	finalHandler := middleware.ErrorHandlingMiddleware(handlerWithAuth)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()

	finalHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var response presenter.CustomError
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ERR_UNAUTHORIZED", response.Code)
	assert.Equal(t, "Missing authorization header", response.Message)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	rec := httptest.NewRecorder()

	handlerWithAuth := middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Adicionar o ErrorHandlingMiddleware
	finalHandler := middleware.ErrorHandlingMiddleware(handlerWithAuth)
	finalHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var response presenter.CustomError
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ERR_UNAUTHORIZED", response.Code)
	assert.Equal(t, "Invalid token", response.Message)
}

func TestAuthMiddleware_InvalidUserID(t *testing.T) {
	claims := &middleware.Claims{
		UserID: 0,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(middleware.JwtKey)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rec := httptest.NewRecorder()

	handlerWithAuth := middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Adicionar o ErrorHandlingMiddleware
	finalHandler := middleware.ErrorHandlingMiddleware(handlerWithAuth)
	finalHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var response presenter.CustomError
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ERR_UNAUTHORIZED", response.Code)
	assert.Equal(t, "Invalid token", response.Message)
}
