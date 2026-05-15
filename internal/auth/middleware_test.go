package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewAuthManager(t *testing.T) {
	am := NewAuthManager()
	if am == nil {
		t.Fatal("Expected non-nil auth manager")
	}
	if am.keys == nil {
		t.Error("Expected keys map to be initialized")
	}
}

func TestAuthManager_GenerateAPIKey(t *testing.T) {
	am := NewAuthManager()
	
	key, err := am.GenerateAPIKey("user123", []string{"read", "write"}, 24*time.Hour)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if key.Key == "" {
		t.Error("Expected non-empty key")
	}
	if key.UserID != "user123" {
		t.Errorf("Expected user123, got %s", key.UserID)
	}
	if len(key.Scopes) != 2 {
		t.Errorf("Expected 2 scopes, got %d", len(key.Scopes))
	}
}

func TestAuthManager_ValidateAPIKey(t *testing.T) {
	am := NewAuthManager()
	key, _ := am.GenerateAPIKey("user123", []string{"read"}, 24*time.Hour)
	
	apiKey, valid := am.ValidateAPIKey(key.Key)
	if !valid {
		t.Error("Expected key to be valid")
	}
	if apiKey == nil {
		t.Error("Expected non-nil API key")
	}
	if apiKey.UserID != "user123" {
		t.Errorf("Expected user123, got %s", apiKey.UserID)
	}
	
	// Test invalid key
	_, valid = am.ValidateAPIKey("invalid-key")
	if valid {
		t.Error("Expected invalid key to be invalid")
	}
}

func TestAuthManager_ValidateExpiredKey(t *testing.T) {
	am := NewAuthManager()
	key, _ := am.GenerateAPIKey("user123", []string{"read"}, -1*time.Hour)
	
	_, valid := am.ValidateAPIKey(key.Key)
	if valid {
		t.Error("Expected expired key to be invalid")
	}
}

func TestAuthManager_RevokeAPIKey(t *testing.T) {
	am := NewAuthManager()
	key, _ := am.GenerateAPIKey("user123", []string{"read"}, 24*time.Hour)
	
	am.RevokeAPIKey(key.Key)
	
	_, valid := am.ValidateAPIKey(key.Key)
	if valid {
		t.Error("Expected revoked key to be invalid")
	}
}

func TestAuthManager_Middleware(t *testing.T) {
	am := NewAuthManager()
	key, _ := am.GenerateAPIKey("user123", []string{"read"}, 24*time.Hour)
	
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	
	middleware := am.Middleware(handler)
	
	// Test valid key
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+key.Key)
	w := httptest.NewRecorder()
	middleware.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	// Test missing header
	req = httptest.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()
	middleware.ServeHTTP(w, req)
	
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
	
	// Test invalid header format
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	w = httptest.NewRecorder()
	middleware.ServeHTTP(w, req)
	
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
	
	// Test invalid key
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer invalid-key")
	w = httptest.NewRecorder()
	middleware.ServeHTTP(w, req)
	
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}
