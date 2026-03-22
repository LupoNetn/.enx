package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/luponetn/.enx/internal/utils"
)

type contextKey string

const UserIDKey contextKey = "user_id"

// AuthMiddleware validates the access token on protected routes
func AuthMiddleware(accessSecret string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.WriteError(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.WriteError(w, http.StatusUnauthorized, "invalid authorization header format, expected: Bearer <token>")
			return
		}

		claims, err := ValidateToken(parts[1], accessSecret)
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		userID, err := StringToUUID(claims.UserID)
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, "invalid token claims")
			return
		}

		// store user id in context so handlers can access it
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserIDFromContext extracts the authenticated user's id from the request context
func GetUserIDFromContext(ctx context.Context) (pgtype.UUID, error) {
	userID, ok := ctx.Value(UserIDKey).(pgtype.UUID)
	if !ok {
		return pgtype.UUID{}, errors.New("user id not found in context")
	}
	return userID, nil
}