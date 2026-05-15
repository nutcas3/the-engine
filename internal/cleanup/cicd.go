package cleanup

import (
	"context"
	"log"
	"os"
	"strings"
)

// GitHubWorkflowRun represents a GitHub Actions workflow run
type GitHubWorkflowRun struct {
	ID         int64  `json:"id"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// GitHubWorkflowRunsResponse represents the API response
type GitHubWorkflowRunsResponse struct {
	TotalCount   int                 `json:"total_count"`
	WorkflowRuns []GitHubWorkflowRun `json:"workflow_runs"`
}

// checkTestsComplete checks if all tests have completed for a test environment
func (cm *CleanupManager) checkTestsComplete(ctx context.Context, environmentName string) bool {
	// Check environment variable for CI/CD system type
	cicdSystem := os.Getenv("CICD_SYSTEM")
	if cicdSystem == "" {
		cicdSystem = "github" // default
	}

	switch strings.ToLower(cicdSystem) {
	case "github":
		return cm.checkGitHubTestsComplete(ctx, environmentName)
	case "jenkins":
		return cm.checkJenkinsTestsComplete(ctx, environmentName)
	case "gitlab":
		return cm.checkGitLabTestsComplete(ctx)
	default:
		log.Printf("Unknown CI/CD system: %s, assuming tests complete", cicdSystem)
		return true
	}
}

// checkNoActivePipelines checks if there are no active CI/CD pipelines
func (cm *CleanupManager) checkNoActivePipelines(ctx context.Context, environmentName string) bool {
	cicdSystem := os.Getenv("CICD_SYSTEM")
	if cicdSystem == "" {
		cicdSystem = "github"
	}

	switch strings.ToLower(cicdSystem) {
	case "github":
		return cm.checkNoActiveGitHubWorkflows(ctx, environmentName)
	case "jenkins":
		return cm.checkNoActiveJenkinsJobs(ctx, environmentName)
	case "gitlab":
		return cm.checkNoActiveGitLabPipelines(ctx, environmentName)
	default:
		return true
	}
}
