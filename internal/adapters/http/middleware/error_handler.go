package middleware

import (
	"encoding/json"
	"github.com/teamcubation/go-items-challenge/errors"
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ErrorHandlingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				err, ok := rec.(*errors.CustomError)
				if !ok {
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"code":    http.StatusInternalServerError,
						"message": "Internal server error",
					})
					return
				}

				w.WriteHeader(mapErrorToStatus(err.StatusCode))
				json.NewEncoder(w).Encode(map[string]interface{}{
					"code":    err.StatusCode,
					"message": err.Message,
					"details": err.Details,
					"time":    err.Timestamp,
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func mapErrorToStatus(code int) int {
	switch code {
	case 400:
		return http.StatusBadRequest
	case 404:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
