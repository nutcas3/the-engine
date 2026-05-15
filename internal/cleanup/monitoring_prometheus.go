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

func (cm *CleanupManager) checkPrometheusIdle(ctx context.Context, env string, threshold time.Duration) bool {
	endpoint := strings.TrimSpace(os.Getenv("PROMETHEUS_URL"))
	if endpoint == "" {
		log.Printf("PROMETHEUS_URL not configured; cannot verify activity for environment %s", env)
		return false
	}

	duration := formatPromDuration(threshold)
	query := fmt.Sprintf(`sum(rate(http_requests_total{environment="%s"}[%s]))`, env, duration)
	queryURL := fmt.Sprintf("%s/api/v1/query?query=%s", strings.TrimRight(endpoint, "/"), url.QueryEscape(query))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, queryURL, nil)
	if err != nil {
		log.Printf("Failed to create Prometheus request: %v", err)
		return false
	}

	if token := strings.TrimSpace(os.Getenv("PROMETHEUS_TOKEN")); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to query Prometheus: %v", err)
		return false
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			Result []struct {
				Value []any `json:"value"`
			} `json:"result"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Failed to decode Prometheus response: %v", err)
		return false
	}

	if len(result.Data.Result) == 0 {
		return true
	}

	for _, series := range result.Data.Result {
		if len(series.Value) < 2 {
			continue
		}
		value, ok := series.Value[1].(string)
		if !ok {
			continue
		}
		if value != "0" && value != "0.0" {
			log.Printf("Prometheus detected activity for environment %s with rate %s", env, value)
			return false
		}
	}

	return true
}

func formatPromDuration(d time.Duration) string {
	seconds := int(d.Seconds())
	if seconds%3600 == 0 {
		return fmt.Sprintf("%dh", seconds/3600)
	}
	if seconds%60 == 0 {
		return fmt.Sprintf("%dm", seconds/60)
	}
	return fmt.Sprintf("%ds", seconds)
}
