package auth

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/kchopper/curlpp/internal/config"
)

// Authenticator interface defines the methods required for different auth types
type Authenticator interface {
	ApplyAuth(*http.Request) error
}

// GetAuthenticator returns the appropriate authenticator based on config
func GetAuthenticator(cfg config.AuthConfig) (Authenticator, error) {
	switch cfg.Type {
	case "basic":
		return &BasicAuth{
			Username: cfg.Username,
			Password: cfg.Password,
		}, nil
	case "bearer":
		return &BearerAuth{
			Token: cfg.Token,
		}, nil
	case "apikey":
		return &APIKeyAuth{
			Token: cfg.Token,
		}, nil
	case "oauth2":
		return &OAuth2Auth{
			Token: cfg.Token,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported auth type: %s", cfg.Type)
	}
}

// BasicAuth implements basic authentication
type BasicAuth struct {
	Username string
	Password string
}

func (b *BasicAuth) ApplyAuth(req *http.Request) error {
	auth := base64.StdEncoding.EncodeToString([]byte(b.Username + ":" + b.Password))
	req.Header.Set("Authorization", "Basic "+auth)
	return nil
}

// BearerAuth implements bearer token authentication
type BearerAuth struct {
	Token string
}

func (b *BearerAuth) ApplyAuth(req *http.Request) error {
	req.Header.Set("Authorization", "Bearer "+b.Token)
	return nil
}

// APIKeyAuth implements API key authentication
type APIKeyAuth struct {
	Token string
}

func (a *APIKeyAuth) ApplyAuth(req *http.Request) error {
	req.Header.Set("X-API-Key", a.Token)
	return nil
}

// OAuth2Auth implements OAuth2 authentication
type OAuth2Auth struct {
	Token        string
	ExpiresAt    time.Time
	RefreshToken string
	TokenURL     string
	ClientID     string
	ClientSecret string
}

func (o *OAuth2Auth) ApplyAuth(req *http.Request) error {
	// Check if token is expired and refresh if needed
	if time.Now().After(o.ExpiresAt) {
		if err := o.refreshToken(); err != nil {
			return fmt.Errorf("failed to refresh token: %w", err)
		}
	}
	req.Header.Set("Authorization", "Bearer "+o.Token)
	return nil
}

func (o *OAuth2Auth) refreshToken() error {
	// Implement OAuth2 token refresh logic
	// This would typically involve making a request to the TokenURL
	// with the refresh token to get a new access token
	return nil
}
