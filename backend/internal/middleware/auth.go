// Package middleware provides HTTP middleware for authentication and authorization.
package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"manpower-backend/internal/ctxkeys"
)

// Auth validates the JWT token from the Authorization header and
// injects the user's ID and role into the request context.
func Auth(jwtSecret string) func(http.Handler) http.Handler {
	secret := []byte(jwtSecret)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from "Authorization: Bearer <token>" header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeError(w, http.StatusUnauthorized, "Authorization header required")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				writeError(w, http.StatusUnauthorized, "Invalid authorization format. Use: Bearer <token>")
				return
			}

			tokenString := parts[1]

			// Parse and validate the JWT
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return secret, nil
			})

			if err != nil || !token.Valid {
				writeError(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			// Extract claims and inject into request context
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				writeError(w, http.StatusUnauthorized, "Invalid token claims")
				return
			}

			userID, _ := claims["userId"].(string)
			role, _ := claims["role"].(string)

			if userID == "" {
				writeError(w, http.StatusUnauthorized, "Invalid token: missing user ID")
				return
			}

			ctx := context.WithValue(r.Context(), ctxkeys.UserID, userID)
			ctx = context.WithValue(ctx, ctxkeys.UserRole, role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// roleLevel maps role names to permission levels for hierarchical access checks.
// Higher numbers mean more permissions.
var roleLevel = map[string]int{
	"viewer": 1,
	"admin":  2,
}

// RequireMinRole returns middleware that restricts access to users with at least
// the specified role level. Role hierarchy: admin > viewer.
// Must be used after Auth middleware (depends on role being in context).
func RequireMinRole(minRole string) func(http.Handler) http.Handler {
	minLevel := roleLevel[minRole]

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, _ := r.Context().Value(ctxkeys.UserRole).(string)
			level := roleLevel[userRole]

			if level < minLevel {
				writeError(w, http.StatusForbidden, "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// writeError writes a JSON error response without importing the handlers package,
// which would create an import cycle.
func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
