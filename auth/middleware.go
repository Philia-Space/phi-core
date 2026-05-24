package auth

import (
	"net/http"
	"strings"
)

// JWTAuthMiddleware implements the Middleware interface using JWT validation.
type JWTAuthMiddleware struct {
	service Service
}

// NewJWTAuthMiddleware creates a new JWT auth middleware.
func NewJWTAuthMiddleware(service Service) *JWTAuthMiddleware {
	return &JWTAuthMiddleware{service: service}
}

// Handler validates the JWT token and stores the user in context.
func (m *JWTAuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if token == "" {
			http.Error(w, "unauthorized: missing token", http.StatusUnauthorized)
			return
		}

		user, err := m.service.Validate(r.Context(), token)
		if err != nil {
			http.Error(w, "unauthorized: invalid token", http.StatusUnauthorized)
			return
		}

		ctx := ContextWithUser(r.Context(), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRoles returns middleware that checks user roles.
func (m *JWTAuthMiddleware) RequireRoles(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := UserFromContext(r.Context())
			if !ok {
				http.Error(w, "unauthorized: no user in context", http.StatusUnauthorized)
				return
			}

			if !hasRole(user.Roles, roles) {
				http.Error(w, "forbidden: insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func extractToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return ""
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	return parts[1]
}

func hasRole(userRoles, requiredRoles []string) bool {
	roleSet := make(map[string]struct{}, len(userRoles))
	for _, r := range userRoles {
		roleSet[r] = struct{}{}
	}
	for _, required := range requiredRoles {
		if _, ok := roleSet[required]; !ok {
			return false
		}
	}
	return true
}

// Ensure JWTAuthMiddleware implements Middleware interface.
var _ Middleware = (*JWTAuthMiddleware)(nil)
