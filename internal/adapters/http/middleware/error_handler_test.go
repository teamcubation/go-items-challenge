package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/teamcubation/go-items-challenge/errors"
	"github.com/teamcubation/go-items-challenge/internal/adapters/http/middleware"
)

func TestErrorHandlingMiddleware_CustomError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic(errors.New(400, "Bad request", map[string]interface{}{
			"field": "username",
		}))
	})

	testHandler := middleware.ErrorHandlingMiddleware(handler)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	testHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(400), response["code"])
	assert.Equal(t, "Bad request", response["message"])
	assert.Equal(t, "username", response["details"].(map[string]interface{})["field"])
}

func TestErrorHandlingMiddleware_NonCustomError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("unexpected error")
	})

	testHandler := middleware.ErrorHandlingMiddleware(handler)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	testHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusInternalServerError), response["code"])
	assert.Equal(t, "Internal server error", response["message"])
}

func TestMapErrorToStatus(t *testing.T) {
	assert.Equal(t, 400, middleware.MapErrorToStatus(400))
	assert.Equal(t, 404, middleware.MapErrorToStatus(404))
	assert.Equal(t, 500, middleware.MapErrorToStatus(500)) // Código padrão
}
