package cleanup

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type argoCDApplication struct {
	Metadata struct {
		Name string `json:"name"`
	} `json:"metadata"`
	Status struct {
		Sync struct {
			Status     string    `json:"status"`
			FinishedAt time.Time `json:"finishedAt"`
		} `json:"sync"`
	} `json:"status"`
}

func (cm *CleanupManager) checkArgoCDDeployments(ctx context.Context, env string, threshold time.Duration) bool {
	argocdURL := strings.TrimSpace(os.Getenv("ARGOCD_URL"))
	argocdToken := strings.TrimSpace(os.Getenv("ARGOCD_TOKEN"))

	if argocdURL == "" || argocdToken == "" {
		log.Printf("ArgoCD credentials not configured; cannot verify deployments for environment %s", env)
		return false
	}

	queryURL := fmt.Sprintf("%s/api/v1/applications?selector=environment=%s", strings.TrimRight(argocdURL, "/"), env)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, queryURL, nil)
	if err != nil {
		log.Printf("Failed to create ArgoCD request: %v", err)
		return false
	}

	req.Header.Set("Authorization", "Bearer "+argocdToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to query ArgoCD: %v", err)
		return false
	}
	defer resp.Body.Close()

	var result struct {
		Items []argoCDApplication `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Failed to decode ArgoCD response: %v", err)
		return false
	}

	cutoff := time.Now().Add(-threshold)
	for _, app := range result.Items {
		if app.Status.Sync.FinishedAt.After(cutoff) {
			log.Printf("ArgoCD application %s synced at %v", app.Metadata.Name, app.Status.Sync.FinishedAt)
			return false
		}
	}

	return true
}
