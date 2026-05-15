package integrations

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// PagerDutyEvent represents a PagerDuty event
type PagerDutyEvent struct {
	RoutingKey  string           `json:"routing_key"`
	EventAction string           `json:"event_action"`
	DedupKey    string           `json:"dedup_key,omitempty"`
	Payload     PagerDutyPayload `json:"payload"`
	Client      string           `json:"client"`
	ClientURL   string           `json:"client_url,omitempty"`
}

// PagerDutyPayload represents the payload of a PagerDuty event
type PagerDutyPayload struct {
	Summary       string         `json:"summary"`
	Severity      string         `json:"severity"`
	Source        string         `json:"source"`
	Timestamp     string         `json:"timestamp"`
	Component     string         `json:"component,omitempty"`
	Group         string         `json:"group,omitempty"`
	CustomDetails map[string]any `json:"custom_details,omitempty"`
}

// PagerDutyClient handles PagerDuty alerts
type PagerDutyClient struct {
	client     *http.Client
	routingKey string
}

// NewPagerDutyClient creates a new PagerDuty client
func NewPagerDutyClient(routingKey string) *PagerDutyClient {
	return &PagerDutyClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		routingKey: routingKey,
	}
}

// Trigger sends a PagerDuty trigger event
func (pdc *PagerDutyClient) Trigger(ctx context.Context, dedupKey, summary, source, component string, severity string, customDetails map[string]any) error {
	event := PagerDutyEvent{
		RoutingKey:  pdc.routingKey,
		EventAction: "trigger",
		DedupKey:    dedupKey,
		Payload: PagerDutyPayload{
			Summary:       summary,
			Severity:      severity,
			Source:        source,
			Timestamp:     time.Now().UTC().Format(time.RFC3339),
			Component:     component,
			CustomDetails: customDetails,
		},
		Client:    "Sovereign Engine",
		ClientURL: "https://github.com/nutcas3/the-engine",
	}

	return pdc.sendEvent(ctx, event)
}

// Resolve sends a PagerDuty resolve event
func (pdc *PagerDutyClient) Resolve(ctx context.Context, dedupKey string) error {
	event := PagerDutyEvent{
		RoutingKey:  pdc.routingKey,
		EventAction: "resolve",
		DedupKey:    dedupKey,
		Payload: PagerDutyPayload{
			Summary:   "Issue resolved",
			Severity:  "info",
			Source:    "Sovereign Engine",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
		Client:    "Sovereign Engine",
		ClientURL: "https://github.com/nutcas3/the-engine",
	}

	return pdc.sendEvent(ctx, event)
}

// sendEvent sends a PagerDuty event
func (pdc *PagerDutyClient) sendEvent(ctx context.Context, event PagerDutyEvent) error {
	jsonData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://events.pagerduty.com/v2/enqueue", bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := pdc.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("PagerDuty returned status %d", resp.StatusCode)
	}

	return nil
}

// TriggerCostAlert triggers a PagerDuty cost alert
func (pdc *PagerDutyClient) TriggerCostAlert(ctx context.Context, environment string, currentCost, threshold float64) error {
	dedupKey := fmt.Sprintf("cost-alert-%s", environment)
	summary := fmt.Sprintf("Cost threshold exceeded in %s environment", environment)

	customDetails := map[string]any{
		"environment":  environment,
		"current_cost": currentCost,
		"threshold":    threshold,
		"overspend":    currentCost - threshold,
	}

	return pdc.Trigger(ctx, dedupKey, summary, "Sovereign Engine", "Budget Monitor", "warning", customDetails)
}

// TriggerCleanupFailure triggers a PagerDuty alert for cleanup failures
func (pdc *PagerDutyClient) TriggerCleanupFailure(ctx context.Context, environment, provider, resourceID string, error string) error {
	dedupKey := fmt.Sprintf("cleanup-failure-%s-%s", environment, resourceID)
	summary := fmt.Sprintf("Cleanup failed for %s in %s environment", resourceID, environment)

	customDetails := map[string]any{
		"environment": environment,
		"provider":    provider,
		"resource_id": resourceID,
		"error":       error,
	}

	return pdc.Trigger(ctx, dedupKey, summary, "Sovereign Engine", "Cleanup Manager", "error", customDetails)
}

// ResolveCostAlert resolves a PagerDuty cost alert
func (pdc *PagerDutyClient) ResolveCostAlert(ctx context.Context, environment string) error {
	dedupKey := fmt.Sprintf("cost-alert-%s", environment)
	return pdc.Resolve(ctx, dedupKey)
}
