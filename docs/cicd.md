# CI/CD Integration Guide

This guide covers configuring CI/CD systems for test completion tracking and active pipeline detection.

## Overview

The cleanup system integrates with CI/CD platforms to:
- Check if tests have completed before cleaning up test environments
- Detect active workflows/jobs/pipelines before nuking environments
- Integrate with GitHub Actions, Jenkins, and GitLab CI

## Supported CI/CD Systems

- **GitHub Actions** - Workflow run status and completion checks
- **Jenkins** - Job status and build completion checks
- **GitLab CI** - Pipeline status and completion checks

## Configuration

### GitHub Actions

#### Environment Variables

```bash
export CICD_SYSTEM=github
export GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxx
export GITHUB_REPOSITORY=owner/repository
```

#### Token Setup

1. Generate a GitHub Personal Access Token:
   - Go to GitHub Settings → Developer settings → Personal access tokens
   - Select scopes: `repo:status`, `repo_deployment`
   - Generate and copy the token

2. Set the environment variable:
   ```bash
   export GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxx
   ```

#### Workflow Naming

The system expects workflows to be named with the environment:
- Test environments: workflows should include environment name in repository
- The system queries all workflow runs for the repository

#### Example Setup

```bash
# Configure GitHub Actions
export CICD_SYSTEM=github
export GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxx
export GITHUB_REPOSITORY=mycompany/myproject

# Test configuration
engine status
```

#### API Endpoints

The system uses the GitHub REST API:
- Get workflow runs: `GET /repos/{owner}/{repo}/actions/runs`
- Check for running/in-progress workflows

### Jenkins

#### Environment Variables

```bash
export CICD_SYSTEM=jenkins
export JENKINS_URL=https://jenkins.example.com
export JENKINS_USER=username
export JENKINS_TOKEN=api_token
```

#### Token Setup

1. Generate a Jenkins API Token:
   - Go to Jenkins → Configure → API Token
   - Generate and copy the token

2. Set the environment variables:
   ```bash
   export JENKINS_USER=myuser
   export JENKINS_TOKEN=11a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2
   export JENKINS_URL=https://jenkins.example.com
   ```

#### Job Naming

The system expects Jenkins jobs to be named with the environment:
- Test environments: jobs should be named `test-{environment}`
- Example: `test-dev`, `test-staging`

#### Example Setup

```bash
# Configure Jenkins
export CICD_SYSTEM=jenkins
export JENKINS_URL=https://jenkins.example.com
export JENKINS_USER=myuser
export JENKINS_TOKEN=11a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2

# Test configuration
engine status
```

#### API Endpoints

The system uses the Jenkins REST API:
- Get job info: `GET /job/{job_name}/api/json`
- Check job status for running builds

### GitLab CI

#### Environment Variables

```bash
export CICD_SYSTEM=gitlab
export GITLAB_URL=https://gitlab.example.com
export GITLAB_TOKEN=glpat-xxxxxxxxxxxxxxxxxxxx
export GITLAB_PROJECT_ID=12345
```

#### Token Setup

1. Generate a GitLab Personal Access Token:
   - Go to GitLab → Settings → Access Tokens
   - Select scopes: `read_api`, `read_repository`
   - Generate and copy the token

2. Set the environment variables:
   ```bash
   export GITLAB_TOKEN=glpat-xxxxxxxxxxxxxxxxxxxx
   export GITLAB_PROJECT_ID=12345
   export GITLAB_URL=https://gitlab.example.com
   ```

#### Project ID

Find your project ID:
- Go to GitLab project → Settings → General
- Project ID is displayed at the top

#### Example Setup

```bash
# Configure GitLab CI
export CICD_SYSTEM=gitlab
export GITLAB_URL=https://gitlab.example.com
export GITLAB_TOKEN=glpat-xxxxxxxxxxxxxxxxxxxx
export GITLAB_PROJECT_ID=12345

# Test configuration
engine status
```

#### API Endpoints

The system uses the GitLab REST API:
- Get pipelines: `GET /projects/{id}/pipelines`
- Check for running/pending pipelines

## How It Works

### Test Completion Check

The `checkTestsComplete` function determines if tests have finished:

```go
testsComplete := cm.checkTestsComplete(ctx, "test-env")
```

This checks for:
- No running workflows (GitHub Actions)
- No building jobs (Jenkins)
- No running/pending pipelines (GitLab CI)

### Active Pipeline Detection

The `checkNoActivePipelines` function detects active CI/CD activity:

```go
noPipelines := cm.checkNoActivePipelines(ctx, "test-env")
```

This checks for:
- No in-progress workflow runs
- No running job builds
- No active pipelines

## Integration with Cleanup Policies

### Test Environments

Test environments check CI/CD status before cleanup:

```yaml
test-auto-cleanup:
  environment: test
  auto_shutdown: true
  shutdown_after: 2h
  nuke_after: 6h
  exclude_patterns: []
```

The system checks:
1. **Tests completed** (via CI/CD integration)
2. **No active pipelines** (via CI/CD integration)
3. **Environment idle for threshold** (via monitoring)

### Dev Environments

Dev environments don't typically use CI/CD checks, but can be configured:

```yaml
dev-auto-cleanup:
  environment: dev
  auto_shutdown: true
  shutdown_after: 8h
  nuke_after: 24h
  exclude_patterns:
    - "essential-*"
```

## Troubleshooting

### GitHub Actions Issues

```bash
# Test GitHub API connectivity
curl -H "Authorization: token $GITHUB_TOKEN" \
  https://api.github.com/repos/$GITHUB_REPOSITORY/actions/runs

# Verify token permissions
curl -H "Authorization: token $GITHUB_TOKEN" \
  https://api.github.com/user
```

### Jenkins Issues

```bash
# Test Jenkins connectivity
curl -u $JENKINS_USER:$JENKINS_TOKEN \
  $JENKINS_URL/api/json

# Check job exists
curl -u $JENKINS_USER:$JENKINS_TOKEN \
  $JENKINS_URL/job/test-dev/api/json
```

### GitLab CI Issues

```bash
# Test GitLab API connectivity
curl -H "PRIVATE-TOKEN: $GITLAB_TOKEN" \
  $GITLAB_URL/api/v4/projects/$GITLAB_PROJECT_ID/pipelines

# Verify token permissions
curl -H "PRIVATE-TOKEN: $GITLAB_TOKEN" \
  $GITLAB_URL/api/v4/user
```

### CI/CD Not Working

If CI/CD checks always return false (tests considered incomplete):

1. Verify environment variables are set
2. Test API connectivity
3. Check for proper job/workflow naming
4. Verify token permissions
5. Check for active pipelines/workflows

## Best Practices

1. **Use descriptive job names** - Include environment name in CI/CD job names
2. **Set appropriate timeouts** - Configure CI/CD jobs with reasonable timeouts
3. **Monitor CI/CD health** - Ensure your CI/CD system is reliable
4. **Use webhook notifications** - Get notified of CI/CD failures
5. **Test with staging** - Test cleanup policies in staging first

## Security Considerations

- Use service accounts with minimal permissions
- Rotate API tokens regularly
- Use read-only tokens where possible
- Store secrets securely (use `engine-encrypt`)
- Monitor for unauthorized API access
- Use IP allowlisting for CI/CD systems

## Examples

### GitHub Actions with Test Environments

```bash
# Configure
export CICD_SYSTEM=github
export GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxx
export GITHUB_REPOSITORY=mycompany/myproject

# The system will:
# 1. Check all workflow runs in the repository
# 2. Detect running/in-progress workflows
# 3. Prevent cleanup if tests are still running
```

### Jenkins with Multi-Environment Jobs

```bash
# Configure
export CICD_SYSTEM=jenkins
export JENKINS_URL=https://jenkins.example.com
export JENKINS_USER=myuser
export JENKINS_TOKEN=11a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2

# The system will:
# 1. Check job named "test-{environment}"
# 2. Detect running builds
# 3. Prevent cleanup if builds are in progress
```

### GitLab CI with Project-Based Pipelines

```bash
# Configure
export CICD_SYSTEM=gitlab
export GITLAB_URL=https://gitlab.example.com
export GITLAB_TOKEN=glpat-xxxxxxxxxxxxxxxxxxxx
export GITLAB_PROJECT_ID=12345

# The system will:
# 1. Check pipelines in the project
# 2. Detect running/pending pipelines
# 3. Prevent cleanup if pipelines are active
```

## Next Steps

- Configure [Deployment Tracking](deployments.md)
- Set up [Monitoring Integration](monitoring.md)
- Define [Cleanup Policies](cleanup.md)
