package integrations

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewSlackClient(t *testing.T) {
	client := NewSlackClient("https://hooks.slack.com/test")
	if client == nil {
		t.Fatal("Expected non-nil client")
	}
	if client.client == nil {
		t.Error("Expected client to have HTTP client")
	}
	if client.webhookURL != "https://hooks.slack.com/test" {
		t.Errorf("Expected https://hooks.slack.com/test, got %s", client.webhookURL)
	}
}

func TestSlackClient_Send(t *testing.T) {
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
	
	client := NewSlackClient(server.URL)
	message := SlackMessage{
		Text: "Test message",
	}
	
	err := client.Send(context.Background(), message)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestSlackClient_SendDeploymentNotification(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	
	client := NewSlackClient(server.URL)
	err := client.SendDeploymentNotification(context.Background(), "test", "aws", "resource-123")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestSlackClient_SendCleanupNotification(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	
	client := NewSlackClient(server.URL)
	err := client.SendCleanupNotification(context.Background(), "test", "aws", "resource-123", 50.0)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestSlackClient_SendCostAlert(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	
	client := NewSlackClient(server.URL)
	err := client.SendCostAlert(context.Background(), "test", 150.0, 100.0)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
