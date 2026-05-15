package integrations

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type WebhookEvent struct {
	EventType   string         `json:"event_type"`
	Timestamp   time.Time      `json:"timestamp"`
	Environment string         `json:"environment"`
	Provider    string         `json:"provider"`
	ResourceID  string         `json:"resource_id"`
	Details     map[string]any `json:"details"`
}

type WebhookClient struct {
	client    *http.Client
	endpoints map[string]string
}

func NewWebhookClient() *WebhookClient {
	return &WebhookClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		endpoints: make(map[string]string),
	}
}

func (wc *WebhookClient) AddEndpoint(name, url string) {
	wc.endpoints[name] = url
}

func (wc *WebhookClient) RemoveEndpoint(name string) {
	wc.endpoints[name] = ""
}

func (wc *WebhookClient) SendEvent(ctx context.Context, event WebhookEvent) error {
	jsonData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	var errors []error
	for name, endpoint := range wc.endpoints {
		if err := wc.sendToEndpoint(ctx, endpoint, jsonData); err != nil {
			errors = append(errors, fmt.Errorf("%s: %w", name, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to send to some endpoints: %v", errors)
	}

	return nil
}

func (wc *WebhookClient) sendToEndpoint(ctx context.Context, endpoint string, data []byte) error {
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Sovereign-Engine/1.0")

	resp, err := wc.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	return nil
}

func CreateDeploymentEvent(environment, provider, resourceID string, details map[string]any) WebhookEvent {
	return WebhookEvent{
		EventType:   "deployment.created",
		Timestamp:   time.Now(),
		Environment: environment,
		Provider:    provider,
		ResourceID:  resourceID,
		Details:     details,
	}
}

func CreateCleanupEvent(environment, provider, resourceID string, details map[string]any) WebhookEvent {
	return WebhookEvent{
		EventType:   "cleanup.completed",
		Timestamp:   time.Now(),
		Environment: environment,
		Provider:    provider,
		ResourceID:  resourceID,
		Details:     details,
	}
}

func CreateCostAlertEvent(environment string, cost float64, threshold float64, details map[string]any) WebhookEvent {
	if details == nil {
		details = make(map[string]any)
	}
	details["cost"] = cost
	details["threshold"] = threshold

	return WebhookEvent{
		EventType:   "cost.alert",
		Timestamp:   time.Now(),
		Environment: environment,
		ResourceID:  "",
		Details:     details,
	}
}
