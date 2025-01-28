package middleware

import (
	"encoding/json"
	"errors"
	"github.com/teamcubation/go-items-challenge/internal/adapters/http/presenter"
	"log"
	"net/http"
)

func ErrorHandlingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic recovered: %v", rec)
				err, ok := rec.(error)
				if !ok {
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"code":    "ERR_INTERNAL_SERVER",
						"message": "Internal server error",
					})
					return
				}
				handleError(w, err)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func handleError(w http.ResponseWriter, err error) {
	var customErr *presenter.CustomError
	if errors.As(err, &customErr) {
		w.WriteHeader(mapErrorToStatus(customErr.Code))
		json.NewEncoder(w).Encode(customErr)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    "ERR_INTERNAL_SERVER",
			"message": "Internal server error",
		})
	}
}

func mapErrorToStatus(code string) int {
	switch code {
	case "ERR_VALIDATION", "ERR_INVALID_REQUEST_BODY":
		return http.StatusBadRequest
	case "ERR_USERNAME_NOT_FOUND", "ERR_INVALID_CREDENTIALS":
		return http.StatusUnauthorized
	case "ERR_USERNAME_EXISTS":
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
