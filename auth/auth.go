package auth

import (
	"context"
	"time"
)

// User represents an authenticated user.
type User struct {
	ID        string    `json:"id"`
	DiscordID string    `json:"discord_id"`
	Name      string    `json:"name"`
	Avatar    string    `json:"avatar"`
	Email     string    `json:"email"`
	Roles     []string  `json:"roles"`
	CreatedAt time.Time `json:"created_at"`
}

// Token represents a JWT token pair.
type Token struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenType    string    `json:"token_type"`
}

// Provider handles OAuth with external providers.
type Provider interface {
	Name() string
	AuthURL(state string) string
	Exchange(ctx context.Context, code string) (*Token, error)
	GetUser(ctx context.Context, token string) (*User, error)
}

// Service handles authentication operations.
type Service interface {
	Login(ctx context.Context, provider string, code string) (*Token, *User, error)
	Refresh(ctx context.Context, refreshToken string) (*Token, error)
	Validate(ctx context.Context, accessToken string) (*User, error)
	Logout(ctx context.Context, userID string) error
}

// Middleware extracts and validates user from request context.
type Middleware interface {
	Handler(next http.Handler) http.Handler
	RequireRoles(roles ...string) func(http.Handler) http.Handler
}

// ContextKey is the key for storing user in request context.
type ContextKey struct{}

// UserFromContext extracts user from context.
func UserFromContext(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(ContextKey{}).(*User)
	return user, ok
}

// ContextWithUser stores user in context.
func ContextWithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, ContextKey{}, user)
}
