package middleware

import (
	"context"
	"github.com/teamcubation/go-items-challenge/internal/adapters/http/presenter"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

type contextKey string

const UserContextKey contextKey = "userID"

var JwtKey = []byte("your_secret_key")

type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// adding authentication logic here

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			ErrorHandlingMiddleware(w, presenter.New("ERR_UNAUTHORIZED", "Missing token", map[string]interface{}{
				"field": "Authorization required",
			}))
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(_ *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})
		if err != nil || !token.Valid {
			ErrorHandlingMiddleware(w, presenter.New("ERR_UNAUTHORIZED", "Invalid token", map[string]interface{}{
				"error": "Token is not valid",
			}))
			return
		}
		if claims.UserID == 0 {
			ErrorHandlingMiddleware(w, presenter.New("ERR_UNAUTHORIZED", "Invalid token", map[string]interface{}{
				"error": "Token is not valid",
			}))
			return
		}
		ctx := context.WithValue(r.Context(), UserContextKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
