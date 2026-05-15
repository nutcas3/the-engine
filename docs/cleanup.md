# Cleanup Policies Guide

This guide covers configuring automated cleanup policies for environment management and cost control.

## Overview

The cleanup system automatically manages environment lifecycle to control costs:
- **Auto-shutdown** - Shut down idle environments after configured thresholds
- **Auto-nuke** - Delete environments after extended inactivity
- **Pattern-based exclusion** - Protect critical resources from cleanup
- **Integration checks** - Verify no active work before cleanup

## Default Policies

### Dev Environment Policy

```yaml
dev-auto-shutdown:
  environment: dev
  auto_shutdown: true
  shutdown_after: 8h
  nuke_after: 24h
  exclude_patterns:
    - "essential-*"
    - "database-*"
```

**Behavior:**
- Shut down after 8 hours of inactivity
- Nuke after 24 hours of inactivity
- Excludes resources matching patterns
- Checks: no recent deployments, no active resources, no user activity

### Test Environment Policy

```yaml
test-auto-cleanup:
  environment: test
  auto_shutdown: true
  shutdown_after: 2h
  nuke_after: 6h
  exclude_patterns: []
```

**Behavior:**
- Shut down after 2 hours
- Nuke after 6 hours
- No exclusions (all test resources can be cleaned)
- Checks: tests completed, no active pipelines, environment idle

## Policy Configuration

### Policy Fields

```yaml
name: string                    # Policy name (required)
environment: string             # Environment type: dev/test/staging/prod (required)
auto_shutdown: boolean          # Enable automatic shutdown (required)
shutdown_after: duration        # Time before shutdown (e.g., 8h, 2h) (required)
nuke_after: duration            # Time before nuke (e.g., 24h, 6h) (optional)
enabled: boolean                # Enable/disable policy (default: true)
exclude_patterns: []string      # Resource patterns to exclude (optional)
```

### Environment Types

- **dev** - Development environments with longer thresholds
- **test** - Test environments with shorter thresholds
- **staging** - Staging environments (never auto-nuked)
- **prod** - Production environments (never auto-nuked)

### Duration Format

Supported duration formats:
- `8h` - 8 hours
- `24h` - 24 hours
- `2h` - 2 hours
- `6h` - 6 hours
- `30m` - 30 minutes
- `1d` - 1 day

## Creating Custom Policies

### Conservative Dev Policy

```yaml
conservative-dev:
  environment: dev
  auto_shutdown: true
  shutdown_after: 12h
  nuke_after: 48h
  exclude_patterns:
    - "essential-*"
    - "database-*"
    - "cache-*"
  enabled: true
```

### Aggressive Test Policy

```yaml
aggressive-test:
  environment: test
  auto_shutdown: true
  shutdown_after: 1h
  nuke_after: 4h
  exclude_patterns: []
  enabled: true
```

### Staging Policy (No Auto-Nuke)

```yaml
staging-manual:
  environment: staging
  auto_shutdown: true
  shutdown_after: 24h
  nuke_after: 0  # Disabled
  exclude_patterns:
    - "essential-*"
  enabled: true
```

### Production Policy (Manual Only)

```yaml
production-manual:
  environment: prod
  auto_shutdown: false
  shutdown_after: 0
  nuke_after: 0
  exclude_patterns:
    - "*"
  enabled: true
```

## Exclusion Patterns

### Pattern Syntax

Exclusion patterns support wildcards:
- `*` - Matches any characters
- `database-*` - Matches resources starting with "database-"
- `essential-*` - Matches resources starting with "essential-"
- `prod-*` - Matches resources starting with "prod-"
- `cache-*` - Matches resources starting with "cache-"

### Common Patterns

```yaml
exclude_patterns:
  - "essential-*"      # Essential infrastructure
  - "database-*"       # Database resources
  - "cache-*"          # Cache/Redis resources
  - "monitoring-*"     # Monitoring resources
  - "logging-*"        # Logging resources
  - "prod-*"           # Production resources
  - "backup-*"         # Backup resources
```

### Pattern Matching Logic

The system checks each resource ID against patterns:
- Exact match: `myresource` matches `myresource`
- Wildcard prefix: `database-main` matches `database-*`
- Case-sensitive: `Database-*` does not match `database-*`

## Integration Checks

### Dev Environment Checks

Before nuking a dev environment, the system checks:

1. **No Recent Deployments** - Via deployment tracking (ArgoCD/Flux/Kubernetes)
2. **No Active Resources** - Via provider resource status checks
3. **No User Activity** - Via monitoring integration (Prometheus/Datadog/CloudWatch)

All checks must pass before nuking.

### Test Environment Checks

Before nuking a test environment, the system checks:

1. **Tests Completed** - Via CI/CD integration (GitHub/Jenkins/GitLab)
2. **No Active Pipelines** - Via CI/CD integration
3. **Environment Idle** - Via monitoring integration

All checks must pass before nuking.

### Staging/Production

Staging and production environments are never auto-nuked by default:
```go
if policy.Environment == EnvironmentStaging || policy.Environment == EnvironmentProd {
    return false  // Never auto-nuke
}
```

## Manual Operations

### Manual Shutdown

```bash
# Manually shutdown a specific resource
curl -X POST http://localhost:8080/api/cleanup/shutdown \
  -d '{"provider":"aws","resource_id":"i-1234567"}'
```

### Manual Nuke

```bash
# Manually nuke an entire environment
curl -X POST http://localhost:8080/api/nuke/environment \
  -d '{"environment":"test","provider":"aws"}'
```

### View Policies

```bash
# View all cleanup policies
curl http://localhost:8080/api/cleanup/policies
```

### View Active Alerts

```bash
# View alerts for a specific environment
curl http://localhost:8080/api/alerts?environment=dev
```

## Cron Jobs

### Automatic Execution

The following cron jobs run automatically:

```yaml
# Hourly: Dev environment cleanup check
schedule: "@hourly"
action: check_and_cleanup
environment: dev

# Hourly: Test environment cleanup check
schedule: "@hourly"
action: check_and_cleanup
environment: test

# Daily: Cost alert verification
schedule: "@daily"
action: verify_cost_alerts
```

### Cron Schedule Formats

- `@hourly` - Every hour
- `@daily` - Every day at midnight
- `@weekly` - Every week
- `@monthly` - Every month
- Custom cron: `*/5 * * * *` (every 5 minutes)

## API Endpoints

### Cleanup Policies

```bash
# Get all policies
GET /api/cleanup/policies

# Get specific policy
GET /api/cleanup/policies/{name}

# Create policy
POST /api/cleanup/policies
Content-Type: application/json
{
  "name": "my-policy",
  "environment": "dev",
  "auto_shutdown": true,
  "shutdown_after": "8h",
  "nuke_after": "24h"
}

# Update policy
PUT /api/cleanup/policies/{name}

# Delete policy
DELETE /api/cleanup/policies/{name}
```

### Manual Operations

```bash
# Manual shutdown
POST /api/cleanup/shutdown
Content-Type: application/json
{
  "provider": "aws",
  "resource_id": "i-1234567"
}

# Manual nuke
POST /api/nuke/environment
Content-Type: application/json
{
  "environment": "test",
  "provider": "aws"
}
```

## Troubleshooting

### Environment Not Cleaning Up

If environments are not being cleaned up:

1. **Check policy is enabled:**
   ```bash
   curl http://localhost:8080/api/cleanup/policies
   ```

2. **Check cron jobs are running:**
   ```bash
   curl http://localhost:8080/api/cron/jobs
   ```

3. **Check integration status:**
   - Monitoring: Verify `MONITORING_SYSTEM` is configured
   - CI/CD: Verify `CICD_SYSTEM` is configured
   - Deployment: Verify `DEPLOYMENT_SYSTEM` is configured

4. **Check integration checks:**
   - Are tests actually completing?
   - Is the environment truly idle?
   - Are there recent deployments?

### Resources Being Excluded

If resources are not being cleaned:

1. **Check exclusion patterns:**
   ```bash
   curl http://localhost:8080/api/cleanup/policies
   ```

2. **Verify resource IDs match patterns:**
   ```bash
   # List resources
   engine list --provider aws
   ```

3. **Adjust patterns as needed**

### Cron Jobs Not Running

If cron jobs are not executing:

1. **Check job status:**
   ```bash
   curl http://localhost:8080/api/cron/jobs
   ```

2. **Check system health:**
   ```bash
   curl http://localhost:8080/api/health/status
   ```

3. **Check logs for errors**

## Best Practices

1. **Start with conservative thresholds** - Use longer thresholds initially
2. **Test in staging first** - Validate policies before production
3. **Use exclusion patterns** - Protect critical infrastructure
4. **Monitor cleanup operations** - Track what's being cleaned
5. **Set up alerts** - Get notified of cleanup actions
6. **Review regularly** - Adjust thresholds based on usage patterns
7. **Document patterns** - Keep track of exclusion patterns
8. **Use manual operations** - For emergency cleanup
9. **Backup before nuking** - Ensure data is safe
10. **Test integrations** - Verify monitoring/CI/CD/deployment checks

## Security Considerations

- **Manual operations require authentication** - Protect manual endpoints
- **Audit cleanup operations** - Log all cleanup actions
- **Use RBAC** - Restrict who can trigger manual cleanup
- **Review exclusion patterns** - Ensure they don't hide security issues
- **Monitor for abuse** - Watch for excessive manual cleanup
- **Secure API endpoints** - Use authentication and rate limiting

## Examples

### Multi-Environment Setup

```yaml
# Dev - Conservative cleanup
dev-policy:
  environment: dev
  auto_shutdown: true
  shutdown_after: 12h
  nuke_after: 48h
  exclude_patterns:
    - "essential-*"
    - "database-*"

# Test - Aggressive cleanup
test-policy:
  environment: test
  auto_shutdown: true
  shutdown_after: 2h
  nuke_after: 6h
  exclude_patterns: []

# Staging - Manual only
staging-policy:
  environment: staging
  auto_shutdown: true
  shutdown_after: 24h
  nuke_after: 0
  exclude_patterns:
    - "essential-*"

# Prod - No automation
prod-policy:
  environment: prod
  auto_shutdown: false
  shutdown_after: 0
  nuke_after: 0
  exclude_patterns:
    - "*"
```

### Team-Specific Policies

```yaml
# Platform team - Longer thresholds
platform-dev:
  environment: dev
  auto_shutdown: true
  shutdown_after: 24h
  nuke_after: 72h
  exclude_patterns:
    - "essential-*"
    - "database-*"
    - "cache-*"

# Feature teams - Shorter thresholds
feature-test:
  environment: test
  auto_shutdown: true
  shutdown_after: 1h
  nuke_after: 4h
  exclude_patterns: []
```

## Next Steps

- Configure [Monitoring Integration](monitoring.md)
- Set up [CI/CD Integration](cicd.md)
- Configure [Deployment Tracking](deployments.md)
- Review [Setup Guide](setup.md)
