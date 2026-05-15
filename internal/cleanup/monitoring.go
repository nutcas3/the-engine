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

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	cwtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

// checkEnvironmentIdle checks if an environment has been idle for the specified duration.
// It supports Prometheus, Datadog, and CloudWatch monitoring backends selected via
// the MONITORING_SYSTEM environment variable.
func (cm *CleanupManager) checkEnvironmentIdle(ctx context.Context, environmentName string, idleThreshold time.Duration) bool {
	system := strings.ToLower(strings.TrimSpace(os.Getenv("MONITORING_SYSTEM")))
	if system == "" {
		system = "prometheus"
	}

	switch system {
	case "prometheus":
		return cm.checkPrometheusIdle(ctx, environmentName, idleThreshold)
	case "datadog":
		return cm.checkDatadogIdle(ctx, environmentName, idleThreshold)
	case "cloudwatch":
		return cm.checkCloudWatchIdle(ctx, environmentName, idleThreshold)
	default:
		log.Printf("Unknown MONITORING_SYSTEM '%s', assuming environment %s is active", system, environmentName)
		return false
	}
}

// checkPrometheusIdle queries Prometheus to determine if any requests were served for the
// environment within the specified threshold window. If no samples are returned or all
// rates are zero, the environment is considered idle.
func (cm *CleanupManager) checkPrometheusIdle(ctx context.Context, env string, threshold time.Duration) bool {
	endpoint := strings.TrimSpace(os.Getenv("PROMETHEUS_URL"))
	if endpoint == "" {
		log.Printf("PROMETHEUS_URL not configured; cannot verify activity for environment %s", env)
		return false
	}

	dur := formatPromDuration(threshold)
	query := fmt.Sprintf(`sum(rate(http_requests_total{environment="%s"}[%s]))`, env, dur)
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
		Status string `json:"status"`
		Data   struct {
			Result []struct {
				Value []interface{} `json:"value"`
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

// checkDatadogIdle queries Datadog events to determine if any events were emitted for the
// environment within the specified threshold. If no events are returned, the environment
// is considered idle.
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
	queryURL := fmt.Sprintf("https://api.%s/api/v1/events?start=%d&end=%d&tags=environment:%s", site, start, now, url.QueryEscape(env))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, queryURL, nil)
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

// checkCloudWatchIdle runs a CloudWatch Logs Insights query against the configured log group
// to determine if any log entries were emitted within the specified threshold window.
func (cm *CleanupManager) checkCloudWatchIdle(ctx context.Context, env string, threshold time.Duration) bool {
	logGroup := strings.TrimSpace(os.Getenv("CLOUDWATCH_LOG_GROUP"))
	if logGroup == "" {
		log.Printf("CLOUDWATCH_LOG_GROUP not configured; cannot verify activity for environment %s", env)
		return false
	}

	cfg, err := awscfg.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("Failed to load AWS configuration: %v", err)
		return false
	}

	client := cloudwatchlogs.NewFromConfig(cfg)
	query := fmt.Sprintf(`fields @timestamp | filter environment="%s" | stats count() as requestCount`, env)
	start := time.Now().Add(-threshold)

	startQueryOutput, err := client.StartQuery(ctx, &cloudwatchlogs.StartQueryInput{
		LogGroupName: aws.String(logGroup),
		StartTime:    aws.Int64(start.Unix()),
		EndTime:      aws.Int64(time.Now().Unix()),
		QueryString:  aws.String(query),
		Limit:        aws.Int32(1),
	})
	if err != nil {
		log.Printf("Failed to start CloudWatch Logs Insights query: %v", err)
		return false
	}

	queryID := startQueryOutput.QueryId
	if queryID == nil {
		log.Printf("CloudWatch Logs Insights did not return a query ID")
		return false
	}

	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return false
		case <-time.After(2 * time.Second):
		}

		results, err := client.GetQueryResults(ctx, &cloudwatchlogs.GetQueryResultsInput{QueryId: queryID})
		if err != nil {
			log.Printf("Failed to fetch CloudWatch Logs Insights results: %v", err)
			return false
		}

		switch results.Status {
		case cwtypes.QueryStatusComplete:
			for _, row := range results.Results {
				for _, field := range row {
					if aws.ToString(field.Field) == "requestCount" {
						if val := aws.ToString(field.Value); val != "" && val != "0" {
							log.Printf("CloudWatch detected %s requests for environment %s", val, env)
							return false
						}
					}
				}
			}
			return true
		case cwtypes.QueryStatusFailed, cwtypes.QueryStatusCancelled:
			log.Printf("CloudWatch Logs Insights query failed with status %s", results.Status)
			return false
		case cwtypes.QueryStatusTimeout:
			log.Printf("CloudWatch Logs Insights query timed out")
			return false
		default:
			// still running; continue polling
		}
	}

	log.Printf("Timed out waiting for CloudWatch Logs Insights results for environment %s", env)
	return false
}

// checkNoUserActivity determines if the environment has seen no user activity within the
// given threshold by delegating to the monitoring system check.
func (cm *CleanupManager) checkNoUserActivity(ctx context.Context, environmentName string, threshold time.Duration) bool {
	return cm.checkEnvironmentIdle(ctx, environmentName, threshold)
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
