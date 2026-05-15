package alerts

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// AlertType represents different types of alerts
type AlertType string

const (
	AlertTypeCost         AlertType = "cost"
	AlertTypeTTL          AlertType = "ttl"
	AlertTypeTestComplete AlertType = "test_complete"
	AlertTypeManual       AlertType = "manual"
)

// AlertSeverity represents the severity level of an alert
type AlertSeverity string

const (
	SeverityInfo     AlertSeverity = "info"
	SeverityWarning  AlertSeverity = "warning"
	SeverityCritical AlertSeverity = "critical"
)

// Alert represents an infrastructure alert
type Alert struct {
	ID          string        `json:"id"`
	Type        AlertType     `json:"type"`
	Severity    AlertSeverity `json:"severity"`
	Environment string        `json:"environment"`
	Resource    string        `json:"resource"`
	Message     string        `json:"message"`
	CreatedAt   time.Time     `json:"created_at"`
	Resolved    bool          `json:"resolved"`
	ResolvedAt  *time.Time    `json:"resolved_at,omitempty"`
}

// AlertManager manages infrastructure alerts
type AlertManager struct {
	alerts       map[string]*Alert
	mu           sync.RWMutex
	alertChan    chan *Alert
	shutdownChan chan struct{}
}

// NewAlertManager creates a new alert manager
func NewAlertManager() *AlertManager {
	return &AlertManager{
		alerts:       make(map[string]*Alert),
		alertChan:    make(chan *Alert, 100),
		shutdownChan: make(chan struct{}),
	}
}

// Start begins the alert manager background processing
func (am *AlertManager) Start(ctx context.Context) {
	go am.processAlerts(ctx)
}

// Stop stops the alert manager
func (am *AlertManager) Stop() {
	close(am.shutdownChan)
}

// processAlerts handles incoming alerts
func (am *AlertManager) processAlerts(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-am.shutdownChan:
			return
		case alert := <-am.alertChan:
			am.mu.Lock()
			am.alerts[alert.ID] = alert
			am.mu.Unlock()

			log.Printf("ALERT [%s] %s: %s (Environment: %s, Resource: %s)",
				alert.Severity, alert.Type, alert.Message, alert.Environment, alert.Resource)
		}
	}
}

// CreateAlert creates a new alert
func (am *AlertManager) CreateAlert(alertType AlertType, severity AlertSeverity, environment, resource, message string) *Alert {
	alert := &Alert{
		ID:          fmt.Sprintf("%d", time.Now().UnixNano()),
		Type:        alertType,
		Severity:    severity,
		Environment: environment,
		Resource:    resource,
		Message:     message,
		CreatedAt:   time.Now(),
		Resolved:    false,
	}

	am.alertChan <- alert
	return alert
}

// GetAlerts returns all alerts
func (am *AlertManager) GetAlerts() []*Alert {
	am.mu.RLock()
	defer am.mu.RUnlock()

	alerts := make([]*Alert, 0, len(am.alerts))
	for _, alert := range am.alerts {
		alerts = append(alerts, alert)
	}
	return alerts
}

// GetAlertsByEnvironment returns alerts for a specific environment
func (am *AlertManager) GetAlertsByEnvironment(environment string) []*Alert {
	am.mu.RLock()
	defer am.mu.RUnlock()

	var alerts []*Alert
	for _, alert := range am.alerts {
		if alert.Environment == environment {
			alerts = append(alerts, alert)
		}
	}
	return alerts
}

// GetUnresolvedAlerts returns unresolved alerts
func (am *AlertManager) GetUnresolvedAlerts() []*Alert {
	am.mu.RLock()
	defer am.mu.RUnlock()

	var alerts []*Alert
	for _, alert := range am.alerts {
		if !alert.Resolved {
			alerts = append(alerts, alert)
		}
	}
	return alerts
}

// ResolveAlert marks an alert as resolved
func (am *AlertManager) ResolveAlert(alertID string) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	alert, exists := am.alerts[alertID]
	if !exists {
		return fmt.Errorf("alert not found")
	}

	now := time.Now()
	alert.Resolved = true
	alert.ResolvedAt = &now

	return nil
}

// AlertRule defines when alerts should be triggered
type AlertRule struct {
	Name        string    `json:"name"`
	Type        AlertType `json:"type"`
	Environment string    `json:"environment"`
	Threshold   any       `json:"threshold"`
	Enabled     bool      `json:"enabled"`
}

// CheckCostAlert checks if cost threshold is exceeded
func (am *AlertManager) CheckCostAlert(environment string, currentSpend, budget float64) {
	if budget > 0 {
		threshold := budget * 0.8 // Alert at 80% of budget
		if currentSpend > threshold {
			severity := SeverityWarning
			if currentSpend > budget*0.9 {
				severity = SeverityCritical
			}

			am.CreateAlert(AlertTypeCost, severity, environment, "budget",
				fmt.Sprintf("Cost alert: Current spend $%.2f exceeds %d%% of budget $%.2f",
					currentSpend, int(threshold/budget*100), budget))
		}
	}
}

// CheckTTLAlert checks if resources have exceeded their TTL
func (am *AlertManager) CheckTTLAlert(environment, resource string, createdAt time.Time, ttl time.Duration) {
	if time.Since(createdAt) > ttl {
		am.CreateAlert(AlertTypeTTL, SeverityWarning, environment, resource,
			fmt.Sprintf("TTL alert: Resource has exceeded TTL of %v", ttl))
	}
}
