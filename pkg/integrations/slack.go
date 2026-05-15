package integrations

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SlackMessage represents a Slack message
type SlackMessage struct {
	Text       string            `json:"text"`
	Username   string            `json:"username,omitempty"`
	IconEmoji  string            `json:"icon_emoji,omitempty"`
	Channel    string            `json:"channel,omitempty"`
	Attachments []SlackAttachment `json:"attachments,omitempty"`
}

// SlackAttachment represents a Slack message attachment
type SlackAttachment struct {
	Color     string            `json:"color"`
	Title     string            `json:"title"`
	Text      string            `json:"text"`
	Fields    []SlackField      `json:"fields,omitempty"`
	Timestamp int64             `json:"ts,omitempty"`
}

// SlackField represents a Slack attachment field
type SlackField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// SlackClient handles Slack notifications
type SlackClient struct {
	client    *http.Client
	webhookURL string
}

// NewSlackClient creates a new Slack client
func NewSlackClient(webhookURL string) *SlackClient {
	return &SlackClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		webhookURL: webhookURL,
	}
}

// Send sends a message to Slack
func (sc *SlackClient) Send(ctx context.Context, message SlackMessage) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", sc.webhookURL, bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := sc.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Slack returned status %d", resp.StatusCode)
	}

	return nil
}

// SendDeploymentNotification sends a deployment notification to Slack
func (sc *SlackClient) SendDeploymentNotification(ctx context.Context, environment, provider, resourceID string) error {
	message := SlackMessage{
		Text:     fmt.Sprintf("Deployment created in %s environment", environment),
		Username: "Sovereign Engine",
		IconEmoji: ":rocket:",
		Attachments: []SlackAttachment{
			{
				Color: "good",
				Title: "Deployment Created",
				Fields: []SlackField{
					{Title: "Environment", Value: environment, Short: true},
					{Title: "Provider", Value: provider, Short: true},
					{Title: "Resource ID", Value: resourceID, Short: false},
				},
				Timestamp: time.Now().Unix(),
			},
		},
	}

	return sc.Send(ctx, message)
}

// SendCleanupNotification sends a cleanup notification to Slack
func (sc *SlackClient) SendCleanupNotification(ctx context.Context, environment, provider, resourceID string, costSaved float64) error {
	message := SlackMessage{
		Text:     fmt.Sprintf("Resource cleaned up in %s environment", environment),
		Username: "Sovereign Engine",
		IconEmoji: ":broom:",
		Attachments: []SlackAttachment{
			{
				Color: "warning",
				Title: "Resource Cleaned Up",
				Fields: []SlackField{
					{Title: "Environment", Value: environment, Short: true},
					{Title: "Provider", Value: provider, Short: true},
					{Title: "Resource ID", Value: resourceID, Short: false},
					{Title: "Cost Saved", Value: fmt.Sprintf("$%.2f", costSaved), Short: true},
				},
				Timestamp: time.Now().Unix(),
			},
		},
	}

	return sc.Send(ctx, message)
}

// SendCostAlert sends a cost alert to Slack
func (sc *SlackClient) SendCostAlert(ctx context.Context, environment string, currentCost, threshold float64) error {
	message := SlackMessage{
		Text:     fmt.Sprintf("Cost alert for %s environment", environment),
		Username: "Sovereign Engine",
		IconEmoji: ":moneybag:",
		Attachments: []SlackAttachment{
			{
				Color: "danger",
				Title: "Budget Threshold Exceeded",
				Fields: []SlackField{
					{Title: "Environment", Value: environment, Short: true},
					{Title: "Current Cost", Value: fmt.Sprintf("$%.2f", currentCost), Short: true},
					{Title: "Threshold", Value: fmt.Sprintf("$%.2f", threshold), Short: true},
					{Title: "Overspend", Value: fmt.Sprintf("$%.2f", currentCost-threshold), Short: true},
				},
				Timestamp: time.Now().Unix(),
			},
		},
	}

	return sc.Send(ctx, message)
}
