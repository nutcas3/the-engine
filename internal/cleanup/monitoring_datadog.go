package cleanup

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func (cm *CleanupManager) checkDatadogIdle(ctx context.Context, env string, threshold time.Duration) bool {
	apiKey := strings.TrimSpace(os.Getenv("DATADOG_API_KEY"))
	appKey := strings.TrimSpace(os.Getenv("DATADOG_APP_KEY"))
	if apiKey == "" || appKey == "" {
		log.Printf("Datadog credentials not configured; cannot verify activity for environment %s", env)
		return false
	}

	site := strings.TrimSpace(os.Getenv("DATADOG_SITE"))
	if site == "" {
		site = "datadoghq.com"
	}

	now := time.Now().Unix()
	start := now - int64(threshold.Seconds())
	queryURL := url.URL{
		Scheme: "https",
		Host:   "api." + site,
		Path:   "/api/v1/events",
	}
	params := queryURL.Query()
	params.Set("start", fmt.Sprintf("%d", start))
	params.Set("end", fmt.Sprintf("%d", now))
	params.Set("tags", "environment:"+env)
	queryURL.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, queryURL.String(), nil)
	if err != nil {
		log.Printf("Failed to create Datadog request: %v", err)
		return false
	}

	req.Header.Set("DD-API-KEY", apiKey)
	req.Header.Set("DD-APPLICATION-KEY", appKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to query Datadog: %v", err)
		return false
	}
	defer resp.Body.Close()

	var result struct {
		Events []json.RawMessage `json:"events"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Failed to decode Datadog response: %v", err)
		return false
	}

	if len(result.Events) == 0 {
		return true
	}

	log.Printf("Datadog reported %d events for environment %s", len(result.Events), env)
	return false
}
