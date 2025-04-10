package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/auth"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/common"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
)

func ExtractAuthToken(header string) (string, error) {
	parts := strings.Split(header, " ")
	if len(parts) != 2 {
		return "", common.ErrorInvalidAuthheaderFormat
	}
	if parts[0] != "Bearer" {
		return "", common.ErrorInvalidAuthheaderFormat
	}
	return parts[1], nil
}

func NewAuthMiddleware(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			token, err := ExtractAuthToken(authHeader)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			userID, err := auth.GetUserIDFromToken(token, secretKey)

			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			if userID == "" {
				http.Error(w, common.ErrorNoUserID.Error(), http.StatusUnauthorized)
				return
			}

			contextWithUser := context.WithValue(r.Context(), UserIDKey, userID)

			// Call the next handler
			next.ServeHTTP(w, r.WithContext(contextWithUser))

		})
	}
}
