package integrations

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewWebhookClient(t *testing.T) {
	client := NewWebhookClient()
	if client == nil {
		t.Fatal("Expected non-nil client")
	}
	if client.client == nil {
		t.Error("Expected client to have HTTP client")
	}
	if client.endpoints == nil {
		t.Error("Expected client to have endpoints map")
	}
}

func TestWebhookClient_AddEndpoint(t *testing.T) {
	client := NewWebhookClient()
	client.AddEndpoint("test", "https://example.com/hook")

	if len(client.endpoints) != 1 {
		t.Errorf("Expected 1 endpoint, got %d", len(client.endpoints))
	}

	url, exists := client.endpoints["test"]
	if !exists {
		t.Error("Expected endpoint to exist")
	}
	if url != "https://example.com/hook" {
		t.Errorf("Expected https://example.com/hook, got %s", url)
	}
}

func TestWebhookClient_RemoveEndpoint(t *testing.T) {
	client := NewWebhookClient()
	client.AddEndpoint("test", "https://example.com/hook")
	client.RemoveEndpoint("test")

	if len(client.endpoints) != 0 {
		t.Errorf("Expected 0 endpoints, got %d", len(client.endpoints))
	}
}

func TestWebhookClient_SendEvent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected application/json content type, got %s", r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewWebhookClient()
	client.AddEndpoint("test", server.URL)

	event := WebhookEvent{
		EventType:   "test.event",
		Timestamp:   time.Now(),
		Environment: "test",
		Provider:    "test",
		ResourceID:  "test-123",
		Details:     map[string]interface{}{"key": "value"},
	}

	err := client.SendEvent(context.Background(), event)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestCreateDeploymentEvent(t *testing.T) {
	details := map[string]interface{}{"tier": "small"}
	event := CreateDeploymentEvent("test", "aws", "resource-123", details)

	if event.EventType != "deployment.created" {
		t.Errorf("Expected deployment.created, got %s", event.EventType)
	}
	if event.Environment != "test" {
		t.Errorf("Expected test, got %s", event.Environment)
	}
	if event.Provider != "aws" {
		t.Errorf("Expected aws, got %s", event.Provider)
	}
	if event.ResourceID != "resource-123" {
		t.Errorf("Expected resource-123, got %s", event.ResourceID)
	}
}

func TestCreateCleanupEvent(t *testing.T) {
	details := map[string]interface{}{"cost_saved": 10.50}
	event := CreateCleanupEvent("test", "aws", "resource-123", details)

	if event.EventType != "cleanup.completed" {
		t.Errorf("Expected cleanup.completed, got %s", event.EventType)
	}
}

func TestCreateCostAlertEvent(t *testing.T) {
	details := map[string]interface{}{}
	event := CreateCostAlertEvent("test", 150.0, 100.0, details)

	if event.EventType != "cost.alert" {
		t.Errorf("Expected cost.alert, got %s", event.EventType)
	}
	if event.Details["cost"] != 150.0 {
		t.Errorf("Expected 150.0, got %v", event.Details["cost"])
	}
	if event.Details["threshold"] != 100.0 {
		t.Errorf("Expected 100.0, got %v", event.Details["threshold"])
	}
}
