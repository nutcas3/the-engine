package cleanup

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	cwtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

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

	startOutput, err := client.StartQuery(ctx, &cloudwatchlogs.StartQueryInput{
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

	queryID := startOutput.QueryId
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
						if count := aws.ToString(field.Value); count != "" && count != "0" {
							log.Printf("CloudWatch detected %s requests for environment %s", count, env)
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
