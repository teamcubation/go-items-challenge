package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/teamcubation/go-items-challenge/internal/adapters/http/presenter"
)

func ErrorHandlingMiddleware(w http.ResponseWriter, err error) {
	var customErr *presenter.CustomError
	switch t := err.(type) {
	case *presenter.CustomError:
		customErr = t
	default:
		customErr = presenter.New("ERR_INTERNAL_SERVER", "Internal server error", nil)
	}

	statusCode := MapErrorToStatus(customErr.Code)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(customErr)
}

func MapErrorToStatus(code string) int {
	switch code {
	case "ERR_UNAUTHORIZED":
		return http.StatusUnauthorized
	case "ERR_VALIDATION", "ERR_INVALID_REQUEST_BODY", "ERR_BAD_REQUEST":
		return http.StatusBadRequest
	case "ERR_USERNAME_NOT_FOUND", "ERR_INVALID_CREDENTIALS":
		return http.StatusUnauthorized
	case "ERR_USERNAME_EXISTS":
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
