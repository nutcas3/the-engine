package cleanup

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// checkGitLabTestsComplete checks GitLab CI for test completion
func (cm *CleanupManager) checkGitLabTestsComplete(ctx context.Context, environmentName string) bool {
	gitlabURL := os.Getenv("GITLAB_URL")
	gitlabToken := os.Getenv("GITLAB_TOKEN")
	projectID := os.Getenv("GITLAB_PROJECT_ID")

	if gitlabURL == "" || gitlabToken == "" || projectID == "" {
		log.Printf("GitLab credentials not configured")
		return false
	}

	url := fmt.Sprintf("%s/api/v4/projects/%s/pipelines?per_page=5", gitlabURL, projectID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Printf("Failed to create GitLab API request: %v", err)
		return false
	}

	req.Header.Set("PRIVATE-TOKEN", gitlabToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to query GitLab API: %v", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("GitLab API returned status %d", resp.StatusCode)
		return false
	}

	var pipelines []struct {
		ID     int    `json:"id"`
		Status string `json:"status"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&pipelines); err != nil {
		log.Printf("Failed to decode GitLab response: %v", err)
		return false
	}

	for _, pipeline := range pipelines {
		if pipeline.Status == "running" || pipeline.Status == "pending" {
			log.Printf("Found active pipeline %d", pipeline.ID)
			return false
		}
	}

	return true
}

// checkNoActiveGitLabPipelines checks for running GitLab pipelines
func (cm *CleanupManager) checkNoActiveGitLabPipelines(ctx context.Context, environmentName string) bool {
	gitlabURL := os.Getenv("GITLAB_URL")
	gitlabToken := os.Getenv("GITLAB_TOKEN")
	projectID := os.Getenv("GITLAB_PROJECT_ID")

	if gitlabURL == "" || gitlabToken == "" {
		return false
	}

	url := fmt.Sprintf("%s/api/v4/projects/%s/pipelines?status=running", gitlabURL, projectID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false
	}

	req.Header.Set("PRIVATE-TOKEN", gitlabToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	var pipelines []interface{}
	if err := json.NewDecoder(resp.Body).Decode(&pipelines); err != nil {
		return false
	}

	return len(pipelines) == 0
}
