package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/teamcubation/go-items-challenge/internal/adapters/http/middleware"
	"github.com/teamcubation/go-items-challenge/internal/adapters/http/presenter"
)

func TestErrorHandlingMiddleware_CustomError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic(presenter.New("ERR_BAD_REQUEST", "Bad request", map[string]interface{}{
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
	assert.Equal(t, "ERR_BAD_REQUEST", response["code"])
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
	assert.Equal(t, "ERR_INTERNAL_SERVER", response["code"])
	assert.Equal(t, "Internal server error", response["message"])
}

func TestMapErrorToStatus(t *testing.T) {
	assert.Equal(t, http.StatusBadRequest, middleware.MapErrorToStatus("ERR_VALIDATION"))
	assert.Equal(t, http.StatusBadRequest, middleware.MapErrorToStatus("ERR_INVALID_REQUEST_BODY"))
	assert.Equal(t, http.StatusUnauthorized, middleware.MapErrorToStatus("ERR_USERNAME_NOT_FOUND"))
	assert.Equal(t, http.StatusConflict, middleware.MapErrorToStatus("ERR_USERNAME_EXISTS"))
	assert.Equal(t, http.StatusInternalServerError, middleware.MapErrorToStatus("UNKNOWN_ERROR_CODE"))
}
