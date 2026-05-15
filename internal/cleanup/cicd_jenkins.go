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

// checkJenkinsTestsComplete checks Jenkins for test completion
func (cm *CleanupManager) checkJenkinsTestsComplete(ctx context.Context, environmentName string) bool {
	jenkinsURL := os.Getenv("JENKINS_URL")
	jenkinsUser := os.Getenv("JENKINS_USER")
	jenkinsToken := os.Getenv("JENKINS_TOKEN")

	if jenkinsURL == "" || jenkinsUser == "" || jenkinsToken == "" {
		log.Printf("Jenkins credentials not configured")
		return false
	}

	jobName := fmt.Sprintf("test-%s", environmentName)
	url := fmt.Sprintf("%s/job/%s/lastBuild/api/json", jenkinsURL, jobName)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Printf("Failed to create Jenkins API request: %v", err)
		return false
	}

	req.SetBasicAuth(jenkinsUser, jenkinsToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to query Jenkins API: %v", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Jenkins API returned status %d", resp.StatusCode)
		return false
	}

	var build struct {
		Building bool   `json:"building"`
		Result   string `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&build); err != nil {
		log.Printf("Failed to decode Jenkins response: %v", err)
		return false
	}

	return !build.Building && (build.Result == "SUCCESS" || build.Result == "UNSTABLE")
}

// checkNoActiveJenkinsJobs checks for running Jenkins jobs
func (cm *CleanupManager) checkNoActiveJenkinsJobs(ctx context.Context, environmentName string) bool {
	jenkinsURL := os.Getenv("JENKINS_URL")
	jenkinsUser := os.Getenv("JENKINS_USER")
	jenkinsToken := os.Getenv("JENKINS_TOKEN")

	if jenkinsURL == "" {
		return false
	}

	url := fmt.Sprintf("%s/api/json?tree=jobs[name,lastBuild[building]]", jenkinsURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false
	}

	req.SetBasicAuth(jenkinsUser, jenkinsToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	var result struct {
		Jobs []struct {
			Name      string `json:"name"`
			LastBuild struct {
				Building bool `json:"building"`
			} `json:"lastBuild"`
		} `json:"jobs"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false
	}

	for _, job := range result.Jobs {
		if job.LastBuild.Building {
			return false
		}
	}

	return true
}
