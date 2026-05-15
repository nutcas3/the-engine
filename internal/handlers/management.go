package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"the-engine/internal/alerts"
)

// HandleAlerts returns alert data
func (h *Handlers) HandleAlerts(w http.ResponseWriter, r *http.Request) {
	environment := r.URL.Query().Get("environment")

	var alerts []*alerts.Alert
	if environment != "" {
		alerts = h.alertManager.GetAlertsByEnvironment(environment)
	} else {
		alerts = h.alertManager.GetAlerts()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alerts)
}

// HandleCleanupPolicies returns cleanup policies
func (h *Handlers) HandleCleanupPolicies(w http.ResponseWriter, r *http.Request) {
	policies := h.cleanupManager.GetPolicies()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(policies)
}

// HandleCronJobs returns cron jobs
func (h *Handlers) HandleCronJobs(w http.ResponseWriter, r *http.Request) {
	jobs := h.cronScheduler.GetJobs()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

// HandleNukeOperations returns nuke operations
func (h *Handlers) HandleNukeOperations(w http.ResponseWriter, r *http.Request) {
	operations := h.nukeManager.GetOperations()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(operations)
}

// HandleManualShutdown manually triggers shutdown of a resource
func (h *Handlers) HandleManualShutdown(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Provider   string `json:"provider"`
		ResourceID string `json:"resource_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.cleanupManager.ManualShutdown(context.Background(), request.Provider, request.ResourceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "shutdown initiated"})
}

// HandleManualNuke manually triggers nuking of an environment
func (h *Handlers) HandleManualNuke(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Environment string `json:"environment"`
		Provider    string `json:"provider"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	operation, err := h.nukeManager.NukeEnvironment(context.Background(), request.Environment, request.Provider)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(operation)
}
