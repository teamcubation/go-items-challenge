package middleware_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/teamcubation/go-items-challenge/internal/adapters/http/middleware"
	"github.com/teamcubation/go-items-challenge/internal/adapters/http/presenter"
)

func TestErrorHandlingMiddleware_CustomError(t *testing.T) {
	rec := httptest.NewRecorder()

	// Simule o uso da função ErrorHandlingMiddleware diretamente.
	err := presenter.New("ERR_BAD_REQUEST", "Bad request", map[string]interface{}{
		"field": "username",
	})
	middleware.ErrorHandlingMiddleware(rec, err)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response map[string]interface{}
	unmarshalErr := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, unmarshalErr)
	assert.Equal(t, "ERR_BAD_REQUEST", response["code"])
	assert.Equal(t, "Bad request", response["message"])
	assert.Equal(t, "username", response["details"].(map[string]interface{})["field"])
}

func TestErrorHandlingMiddleware_NonCustomError(t *testing.T) {
	rec := httptest.NewRecorder()

	// Aqui passamos um erro genérico para ver como o middleware lida com isso.
	err := errors.New("unexpected error")
	middleware.ErrorHandlingMiddleware(rec, err)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var response map[string]interface{}
	unmarshalErr := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, unmarshalErr)
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
