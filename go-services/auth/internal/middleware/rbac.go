package middleware

import (
	"context"
	"net/http"

	"github.com/ahmetsah/industrial-historian/go-services/auth/internal/service"
	"github.com/golang-jwt/jwt/v5"
)

type ContextKey string

const (
	UserRoleKey ContextKey = "userRole"
	UserIDKey   ContextKey = "userID"
)

type RBACMiddleware struct {
	TokenService *service.TokenService
}

func NewRBACMiddleware(tokenService *service.TokenService) *RBACMiddleware {
	return &RBACMiddleware{TokenService: tokenService}
}

func (m *RBACMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}
		tokenString := authHeader[7:]

		token, err := m.TokenService.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserRoleKey, claims["role"])
		ctx = context.WithValue(ctx, UserIDKey, claims["sub"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *RBACMiddleware) RequireRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(UserRoleKey).(string)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Allow ADMIN to access everything
			if role == "ADMIN" {
				next.ServeHTTP(w, r)
				return
			}

			if role != requiredRole {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
