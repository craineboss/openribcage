// Package auth provides authentication and authorization for A2A communication.
//
// This package handles various authentication schemes used in A2A protocol
// communication including Bearer tokens, API keys, and OAuth2 flows.
package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

// AuthType represents different authentication types
type AuthType string

const (
	AuthTypeNone   AuthType = "none"
	AuthTypeBearer AuthType = "bearer"
	AuthTypeAPIKey AuthType = "apikey"
	AuthTypeOAuth2 AuthType = "oauth2"
)

// Authenticator handles A2A authentication
type Authenticator struct {
	logger *logrus.Logger
}

// NewAuthenticator creates a new A2A authenticator
func NewAuthenticator() *Authenticator {
	return &Authenticator{
		logger: logrus.New(),
	}
}

// Credentials holds authentication credentials
type Credentials struct {
	Type     AuthType          `json:"type"`
	Token    string            `json:"token,omitempty"`
	APIKey   string            `json:"api_key,omitempty"`
	Username string            `json:"username,omitempty"`
	Password string            `json:"password,omitempty"`
	Headers  map[string]string `json:"headers,omitempty"`
}

// AddAuthHeaders adds authentication headers to an HTTP request
func (a *Authenticator) AddAuthHeaders(req *http.Request, creds *Credentials) error {
	if creds == nil {
		return nil
	}

	a.logger.Debugf("Adding authentication headers for type: %s", creds.Type)

	switch creds.Type {
	case AuthTypeNone:
		// No authentication required
		return nil

	case AuthTypeBearer:
		if creds.Token == "" {
			return fmt.Errorf("bearer token is required")
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", creds.Token))

	case AuthTypeAPIKey:
		if creds.APIKey == "" {
			return fmt.Errorf("API key is required")
		}
		// Common API key header names
		req.Header.Set("X-API-Key", creds.APIKey)
		req.Header.Set("Authorization", fmt.Sprintf("ApiKey %s", creds.APIKey))

	case AuthTypeOAuth2:
		// TODO: Implement OAuth2 flow
		// This will be implemented when needed
		return fmt.Errorf("OAuth2 authentication not yet implemented")

	default:
		return fmt.Errorf("unsupported authentication type: %s", creds.Type)
	}

	// Add any custom headers
	for key, value := range creds.Headers {
		req.Header.Set(key, value)
	}

	return nil
}

// ValidateCredentials validates authentication credentials
func (a *Authenticator) ValidateCredentials(creds *Credentials) error {
	if creds == nil {
		return fmt.Errorf("credentials cannot be nil")
	}

	switch creds.Type {
	case AuthTypeNone:
		return nil

	case AuthTypeBearer:
		if strings.TrimSpace(creds.Token) == "" {
			return fmt.Errorf("bearer token cannot be empty")
		}

	case AuthTypeAPIKey:
		if strings.TrimSpace(creds.APIKey) == "" {
			return fmt.Errorf("API key cannot be empty")
		}

	case AuthTypeOAuth2:
		return fmt.Errorf("OAuth2 validation not yet implemented")

	default:
		return fmt.Errorf("unsupported authentication type: %s", creds.Type)
	}

	return nil
}

// LoadCredentialsFromEnv loads credentials from environment variables
func (a *Authenticator) LoadCredentialsFromEnv(prefix string) *Credentials {
	// TODO: Implement environment variable credential loading
	// This will be implemented when needed
	// Look for variables like:
	// - {PREFIX}_AUTH_TYPE
	// - {PREFIX}_TOKEN
	// - {PREFIX}_API_KEY
	// etc.

	return &Credentials{
		Type: AuthTypeNone,
	}
}

// RefreshToken refreshes an authentication token if needed
func (a *Authenticator) RefreshToken(ctx context.Context, creds *Credentials) error {
	// TODO: Implement token refresh logic
	// This will be implemented when needed for OAuth2 flows
	return fmt.Errorf("token refresh not yet implemented")
}
