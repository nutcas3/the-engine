package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"
	"sync"
	"time"
)

// APIKey represents an API key
type APIKey struct {
	Key        string
	UserID     string
	CreatedAt  time.Time
	ExpiresAt  time.Time
	Scopes     []string
}

// AuthManager manages authentication
type AuthManager struct {
	keys map[string]*APIKey
	mu   sync.RWMutex
}

// NewAuthManager creates a new auth manager
func NewAuthManager() *AuthManager {
	return &AuthManager{
		keys: make(map[string]*APIKey),
	}
}

// GenerateAPIKey generates a new API key
func (a *AuthManager) GenerateAPIKey(userID string, scopes []string, ttl time.Duration) (*APIKey, error) {
	keyBytes := make([]byte, 32)
	if _, err := rand.Read(keyBytes); err != nil {
		return nil, err
	}

	key := base64.URLEncoding.EncodeToString(keyBytes)

	apiKey := &APIKey{
		Key:       key,
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(ttl),
		Scopes:    scopes,
	}

	a.mu.Lock()
	a.keys[key] = apiKey
	a.mu.Unlock()

	return apiKey, nil
}

// ValidateAPIKey validates an API key
func (a *AuthManager) ValidateAPIKey(key string) (*APIKey, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	apiKey, found := a.keys[key]
	if !found {
		return nil, false
	}

	if time.Now().After(apiKey.ExpiresAt) {
		return nil, false
	}

	return apiKey, true
}

// RevokeAPIKey revokes an API key
func (a *AuthManager) RevokeAPIKey(key string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.keys, key)
}

// Middleware is authentication middleware
func (a *AuthManager) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		_, valid := a.ValidateAPIKey(token)
		if !valid {
			http.Error(w, "Invalid API key", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
