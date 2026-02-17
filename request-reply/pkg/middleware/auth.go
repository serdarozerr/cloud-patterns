package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/serdarozerr/request-reply/internal/config"
)

type Middleware func(http.Handler) http.Handler

func extractBearerToken(r *http.Request) (string, error) {

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is missing")
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", fmt.Errorf("invalid authorization header format")
	}
	return parts[1], nil
}

func AuthMiddleware(c *config.Config) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessToken, err := extractBearerToken(r)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// In a real application, you would verify the access token here (e.g., by checking it against a database or an external authentication service).
			ctx := context.WithValue(r.Context(), "user", accessToken)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
