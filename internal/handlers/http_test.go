package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"the-engine/internal/kubernetes"
)

func TestHandleIndex(t *testing.T) {
	k8sClient := kubernetes.NewClientOrMock()
	h := NewHandlers(k8sClient)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	h.HandleIndex(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "text/html" {
		t.Errorf("Expected text/html, got %s", contentType)
	}
}

func TestHandleHealth(t *testing.T) {
	k8sClient := kubernetes.NewClientOrMock()
	h := NewHandlers(k8sClient)

	req := httptest.NewRequest("GET", "/api/health/status", nil)
	w := httptest.NewRecorder()

	h.HandleHealth(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected application/json, got %s", contentType)
	}
}

func TestHandleSwagger(t *testing.T) {
	k8sClient := kubernetes.NewClientOrMock()
	h := NewHandlers(k8sClient)

	req := httptest.NewRequest("GET", "/api/swagger", nil)
	w := httptest.NewRecorder()

	h.HandleSwagger(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected application/json, got %s", contentType)
	}
}

func TestHandleCompositions(t *testing.T) {
	// Skip test if compositions directory doesn't exist
	if _, err := os.Stat("compositions"); os.IsNotExist(err) {
		t.Skip("compositions directory not found, skipping test")
	}

	k8sClient := kubernetes.NewClientOrMock()
	h := NewHandlers(k8sClient)

	req := httptest.NewRequest("GET", "/api/compositions", nil)
	w := httptest.NewRecorder()

	h.HandleCompositions(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected application/json, got %s", contentType)
	}
}
