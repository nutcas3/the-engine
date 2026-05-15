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

// checkGitHubTestsComplete checks GitHub Actions for test completion
func (cm *CleanupManager) checkGitHubTestsComplete(ctx context.Context, environmentName string) bool {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Printf("GITHUB_TOKEN not set, cannot check test status")
		return false
	}

	repo := os.Getenv("GITHUB_REPOSITORY")
	if repo == "" {
		log.Printf("GITHUB_REPOSITORY not set")
		return false
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/actions/runs?per_page=10&status=completed", repo)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Printf("Failed to create GitHub API request: %v", err)
		return false
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to query GitHub API: %v", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("GitHub API returned status %d", resp.StatusCode)
		return false
	}

	var result GitHubWorkflowRunsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Failed to decode GitHub API response: %v", err)
		return false
	}

	for _, run := range result.WorkflowRuns {
		if run.Status != "completed" {
			log.Printf("Found incomplete workflow run %d for environment %s", run.ID, environmentName)
			return false
		}
		if run.Conclusion == "failure" || run.Conclusion == "cancelled" {
			log.Printf("Workflow run %d failed or was cancelled", run.ID)
		}
	}

	return true
}

// checkNoActiveGitHubWorkflows checks for running GitHub workflows
func (cm *CleanupManager) checkNoActiveGitHubWorkflows(ctx context.Context, environmentName string) bool {
	token := os.Getenv("GITHUB_TOKEN")
	repo := os.Getenv("GITHUB_REPOSITORY")

	if token == "" || repo == "" {
		return false
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/actions/runs?status=in_progress", repo)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	var result GitHubWorkflowRunsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false
	}

	return result.TotalCount == 0
}
