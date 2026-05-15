# Monitoring Integration Guide

This guide covers configuring monitoring backends for environment idle detection and user activity tracking.

## Overview

The cleanup system uses monitoring backends to determine if environments are idle before initiating cleanup operations. This prevents accidental cleanup of active environments.

## Supported Monitoring Backends

- **Prometheus** - Query metrics for API request rates and activity
- **Datadog** - Check events and logs for recent activity
- **CloudWatch** - Query CloudWatch Logs Insights for log entries

## Configuration

### Prometheus

#### Environment Variables

```bash
export MONITORING_SYSTEM=prometheus
export PROMETHEUS_URL=https://prometheus.example.com
export PROMETHEUS_TOKEN=your_token
```

#### Query Configuration

The system queries Prometheus for HTTP request rates:

```promql
sum(rate(http_requests_total{environment="<env>"}[<duration>]))
```

#### Example Setup

```bash
# Configure Prometheus
export MONITORING_SYSTEM=prometheus
export PROMETHEUS_URL=https://prometheus.monitoring.svc.cluster.local
export PROMETHEUS_TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

# Test configuration
engine status
```

#### Prometheus Query Requirements

Your Prometheus instance must collect:
- `http_requests_total` metric with `environment` label
- Request rate data for the threshold period

### Datadog

#### Environment Variables

```bash
export MONITORING_SYSTEM=datadog
export DATADOG_API_KEY=your_api_key
export DATADOG_APP_KEY=your_app_key
export DATADOG_SITE=datadoghq.com
```

#### Example Setup

```bash
# Configure Datadog
export MONITORING_SYSTEM=datadog
export DATADOG_API_KEY=abcd1234efgh5678
export DATADOG_APP_KEY=ijkl9012mnop3456
export DATADOG_SITE=datadoghq.com

# Test configuration
engine status
```

#### Datadog Event Requirements

The system queries Datadog for events tagged with the environment name. Ensure your applications emit events with:
- `environment:<env>` tag
- Relevant activity events (API calls, user actions, etc.)

### CloudWatch

#### Environment Variables

```bash
export MONITORING_SYSTEM=cloudwatch
export CLOUDWATCH_LOG_GROUP=/aws/lambda/your-app
```

#### AWS Configuration

The system uses AWS SDK v2 with default credential chain:
- Environment variables (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`)
- Shared credentials file (`~/.aws/credentials`)
- IAM role (when running in AWS)
- EC2 instance profile

#### Example Setup

```bash
# Configure CloudWatch
export MONITORING_SYSTEM=cloudwatch
export CLOUDWATCH_LOG_GROUP=/aws/lambda/production

# Configure AWS credentials (if not using IAM role)
export AWS_ACCESS_KEY_ID=AKIA...
export AWS_SECRET_ACCESS_KEY=secret...
export AWS_REGION=us-east-1

# Test configuration
engine status
```

#### CloudWatch Logs Requirements

Your log group must contain:
- Log entries with `environment` field
- Application logs showing API requests and user activity
- Sufficient retention period for threshold queries

## How It Works

### Environment Idle Detection

The `checkEnvironmentIdle` function determines if an environment has been idle for a specified threshold:

```go
idle := cm.checkEnvironmentIdle(ctx, "dev", 8*time.Hour)
```

This checks for:
- No API requests in the threshold period
- No log entries in the threshold period
- No user activity in the threshold period

### User Activity Detection

The `checkNoUserActivity` function checks for user activity before nuking dev environments:

```go
noActivity := cm.checkNoUserActivity(ctx, "dev", 24*time.Hour)
```

This uses the same monitoring backend as idle detection.

## Integration with Cleanup Policies

### Dev Environments

Dev environments check for user activity before nuking:

```yaml
dev-auto-cleanup:
  environment: dev
  auto_shutdown: true
  shutdown_after: 8h
  nuke_after: 24h
  exclude_patterns:
    - "essential-*"
```

The system checks:
1. No deployments in last 24 hours
2. No active resources
3. **No user activity in last 24 hours** (via monitoring)

### Test Environments

Test environments check for activity before cleanup:

```yaml
test-auto-cleanup:
  environment: test
  auto_shutdown: true
  shutdown_after: 2h
  nuke_after: 6h
  exclude_patterns: []
```

The system checks:
1. Tests completed
2. No active pipelines
3. **Environment idle for threshold** (via monitoring)

## Troubleshooting

### Prometheus Issues

```bash
# Check Prometheus connectivity
curl -H "Authorization: Bearer $PROMETHEUS_TOKEN" \
  https://prometheus.example.com/api/v1/query?query=up

# Verify metric exists
curl -H "Authorization: Bearer $PROMETHEUS_TOKEN" \
  "https://prometheus.example.com/api/v1/query?query=http_requests_total"
```

### Datadog Issues

```bash
# Test Datadog API key
curl -X GET "https://api.datadoghq.com/api/v1/validate" \
  -H "DD-API-KEY: $DATADOG_API_KEY" \
  -H "DD-APPLICATION-KEY: $DATADOG_APP_KEY"

# Check for recent events
curl -X GET "https://api.datadoghq.com/api/v1/events" \
  -H "DD-API-KEY: $DATADOG_API_KEY" \
  -H "DD-APPLICATION-KEY: $DATADOG_APP_KEY" \
  -d "tags=environment:dev"
```

### CloudWatch Issues

```bash
# Verify AWS credentials
aws sts get-caller-identity

# Check log group exists
aws logs describe-log-groups --log-group-name-prefix /aws/lambda

# Test query
aws logs start-query \
  --log-group-name /aws/lambda/your-app \
  --start-time $(date -d '1 hour ago' +%s) \
  --end-time $(date +%s) \
  --query-string 'fields @timestamp | limit 5'
```

### Monitoring Not Working

If monitoring checks always return false (environment considered active):

1. Verify environment variables are set
2. Test API connectivity
3. Check for required metrics/logs/events
4. Verify query syntax
5. Check authentication/permissions

## Best Practices

1. **Use appropriate thresholds** - Set thresholds based on your team's workflow
2. **Monitor the monitoring** - Ensure your monitoring system is reliable
3. **Test with staging** - Test cleanup policies in staging first
4. **Use exclusion patterns** - Protect critical resources
5. **Alert on failures** - Get notified if monitoring checks fail

## Security Considerations

- Use service accounts with minimal permissions
- Rotate API keys regularly
- Use IAM roles in AWS instead of access keys
- Store secrets securely (use `engine-encrypt`)
- Monitor for unauthorized API access

## Next Steps

- Configure [CI/CD Integration](cicd.md)
- Set up [Deployment Tracking](deployments.md)
- Define [Cleanup Policies](cleanup.md)
