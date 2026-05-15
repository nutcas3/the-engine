# Sovereign Engine: Multi-Cloud Infrastructure Platform

A minimalist, high-performance Multi-Cloud Infrastructure Platform that abstracts cloud complexity across AWS, Azure, GCP, Hetzner, OVH, and DigitalOcean while maintaining strict governance and cost control.

## Quick Start

### Installation (No Go/Kubernetes Required)

```bash
# Download and install pre-built binaries
curl -fsSL https://raw.githubusercontent.com/nutcas3/the-engine/main/install.sh | bash

# Set your master key (generated during install)
export ENGINE_MASTER_KEY=<your_generated_key>

# Encrypt your cloud credentials
engine-encrypt -key AWS_ACCESS_KEY_ID -value 'your_key' -file ~/.engine/secure.env
engine-encrypt -key AWS_SECRET_ACCESS_KEY -value 'your_secret' -file ~/.engine/secure.env

# Deploy infrastructure
engine deploy --provider aws --tier micro --region us-east-1

# Or use the web interface
engine-web
# Open http://localhost:8080
```

## ⚠️ Required Setup: Alerts & Automated Cleanup

**IMPORTANT**: The Engine includes automated cleanup and alerting systems that **MUST** be configured before deploying to dev/test environments.

### Alert System Configuration

The alert system monitors cost thresholds, resource TTLs, and test completion status:

```go
// Default alert thresholds (customize in your deployment)
- Cost alerts trigger at 80% of budget (warning) and 90% (critical)
- TTL alerts trigger when resources exceed configured lifetime
- Test completion alerts notify when tests finish
```

### Cleanup Policies (REQUIRED)

**Dev environments** are automatically cleaned up after 8 hours of inactivity and nuked after 24 hours.
**Test environments** are automatically cleaned up after 2 hours and nuked after 6 hours.

Configure cleanup policies in your deployment:

```yaml
# Default cleanup policies
dev-auto-shutdown:
  environment: dev
  auto_shutdown: true
  shutdown_after: 8h
  nuke_after: 24h
  exclude_patterns:
    - "essential-*"
    - "database-*"

test-auto-cleanup:
  environment: test
  auto_shutdown: true
  shutdown_after: 2h
  nuke_after: 6h
  exclude_patterns: []
```

### Cron Jobs (Auto-configured)

The following cron jobs run automatically:
- **Hourly**: Dev environment cleanup check
- **Hourly**: Test environment cleanup check
- **Daily**: Cost alert verification

### Production Integrations

The cleanup system integrates with external monitoring, CI/CD, and deployment systems for production-ready environment management.

#### Monitoring Integration

Configure monitoring backend via `MONITORING_SYSTEM` environment variable:

```bash
# Prometheus (default)
export MONITORING_SYSTEM=prometheus
export PROMETHEUS_URL=https://prometheus.example.com
export PROMETHEUS_TOKEN=your_token

# Datadog
export MONITORING_SYSTEM=datadog
export DATADOG_API_KEY=your_api_key
export DATADOG_APP_KEY=your_app_key
export DATADOG_SITE=datadoghq.com

# CloudWatch
export MONITORING_SYSTEM=cloudwatch
export CLOUDWATCH_LOG_GROUP=/aws/lambda/your-app
```

The monitoring integration checks for:
- Environment idle time (no API requests, logs, or user activity)
- Recent user activity before nuking dev environments
- Test environment activity before cleanup

#### CI/CD Integration

Configure CI/CD system via `CICD_SYSTEM` environment variable:

```bash
# GitHub Actions
export CICD_SYSTEM=github
export GITHUB_TOKEN=ghp_xxxxx
export GITHUB_REPOSITORY=owner/repo

# Jenkins
export CICD_SYSTEM=jenkins
export JENKINS_URL=https://jenkins.example.com
export JENKINS_USER=username
export JENKINS_TOKEN=api_token

# GitLab CI
export CICD_SYSTEM=gitlab
export GITLAB_URL=https://gitlab.example.com
export GITLAB_TOKEN=glpat_xxxxx
export GITLAB_PROJECT_ID=12345
```

The CI/CD integration checks for:
- Test completion before cleaning up test environments
- Active workflows/jobs/pipelines before nuking
- Integration with GitHub Actions, Jenkins, and GitLab CI APIs

#### Deployment Tracking

Configure deployment system via `DEPLOYMENT_SYSTEM` environment variable:

```bash
# ArgoCD
export DEPLOYMENT_SYSTEM=argocd
export ARGOCD_URL=https://argocd.example.com
export ARGOCD_TOKEN=your_token

# Flux
export DEPLOYMENT_SYSTEM=flux
# Uses Kubernetes client and Flux CRDs

# Kubernetes
export DEPLOYMENT_SYSTEM=kubernetes
# Uses Kubernetes client for native deployments
```

The deployment tracking checks for:
- Recent deployments before nuking dev environments
- Sync status and timestamps for ArgoCD applications
- Flux Kustomization and HelmRelease reconciliation times
- Kubernetes deployment update timestamps

#### Kubernetes Client Configuration

For Flux and Kubernetes deployment tracking, configure Kubernetes access:

```bash
# In-cluster configuration (default)
# Uses service account when running in Kubernetes

# Local development
export KUBECONFIG=/path/to/kubeconfig
```

### Resource Exclusion Patterns

Protect critical resources from automatic cleanup using wildcard patterns:
- `essential-*` - Excludes all resources starting with "essential-"
- `database-*` - Excludes all database resources
- `prod-*` - Excludes production resources

### Manual Operations

```bash
# Manually shutdown a resource
curl -X POST http://localhost:8080/api/cleanup/shutdown \
  -d '{"provider":"aws","resource_id":"i-1234567"}'

# Manually nuke an environment
curl -X POST http://localhost:8080/api/nuke/environment \
  -d '{"environment":"test","provider":"aws"}'

# View active alerts
curl http://localhost:8080/api/alerts?environment=dev

# View cleanup policies
curl http://localhost:8080/api/cleanup/policies

# View scheduled jobs
curl http://localhost:8080/api/cron/jobs
```

### CLI Usage

```bash
# Deploy new instance
engine deploy --provider aws --tier micro --region us-east-1 --budget 100 --team platform

# List deployments
engine list --provider aws

# Check cost report
engine cost --team platform --month 2024-01

# Check health status
engine status

# Get help
engine help
```

### Web Interface

```bash
# Start web server
engine-web

# Access at http://localhost:8080
# Features:
# - Real-time deployment monitoring
# - Budget tracking and cost visualization
# - Composition management
# - Health status dashboard
# - API documentation (Swagger)
# - Dark/Light theme toggle
```

## Features

### Core Capabilities
- **Multi-Cloud Abstraction**: Single API for 6 cloud providers
- **CLI-First Design**: Standalone binary, no Kubernetes required
- **Web Interface**: Modern dashboard with HTMX + Tailwind v4
- **Secure Environment Management**: AES-GCM encryption for credentials
- **Real-Time Updates**: SSE streaming for live data
- **Rate Limiting**: Token bucket rate limiting for API protection
- **Connection Pooling**: Optimized HTTP connections
- **Comprehensive Health Checks**: Component-level monitoring

### Security Features
- **AES-GCM Encryption**: Secure credential storage
- **Master Key System**: Secure key derivation
- **File Permissions**: Secure file handling (600 permissions)
- **Rate Limiting**: API protection against abuse

## Installation

### Prerequisites
- curl or wget (for downloading binaries)
- No Go, Kubernetes, or Docker required for CLI usage

### Manual Installation

```bash
# Download for Linux
wget https://github.com/nutcas3/the-engine/releases/download/v1.0.0/engine-linux-amd64.tar.gz
tar -xzf engine-linux-amd64.tar.gz
sudo cp engine /usr/local/bin/
sudo cp engine-encrypt /usr/local/bin/
sudo cp engine-web /usr/local/bin/

# Download for macOS
wget https://github.com/nutcas3/the-engine/releases/download/v1.0.0/engine-darwin-amd64.tar.gz
tar -xzf engine-darwin-amd64.tar.gz
sudo cp engine /usr/local/bin/
sudo cp engine-encrypt /usr/local/bin/
sudo cp engine-web /usr/local/bin/
```

### Configuration

Configuration is stored in `~/.engine/config.yaml`:

```yaml
version: 1.0.0
compositions_dir: ./compositions
providers:
  - aws
  - azure
  - gcp
  - hetzner
  - ovh
  - digitalocean
secure_env_file: ~/.engine/secure.env
web_port: 8080
```

### Secure Environment Variables

```bash
# Generate master key
openssl rand -base64 32
export ENGINE_MASTER_KEY=<generated_key>

# Encrypt credentials
engine-encrypt -key AWS_ACCESS_KEY_ID -value 'your_key' -file ~/.engine/secure.env
engine-encrypt -key AWS_SECRET_ACCESS_KEY -value 'your_secret' -file ~/.engine/secure.env

# The encrypted values are stored in ~/.engine/secure.env
# They are automatically decrypted when needed
```

## Architecture

### CLI-First Approach
```
User (CLI Tool)
    |
    v
Secure Environment (AES-GCM Encrypted)
    |
    v
Cloud Provider APIs
    |
    v
Infrastructure Provisioning
```

### Web Interface
```
User (Browser)
    |
    v
Web Server (engine-web)
    |
    v
HTMX + Tailwind v4 UI
    |
    v
API Endpoints (Rate Limited)
    |
    v
Secure Environment (AES-GCM Encrypted)
    |
    v
Cloud Provider APIs
```

## Provider Support

| Provider | Status | Tiers Supported |
|----------|--------|-----------------|
| AWS | Complete | micro, small, pro |
| Azure | In Progress | micro, small, pro |
| GCP | In Progress | micro, small, pro |
| Hetzner | Complete | micro, small, pro |
| OVH | In Progress | micro, small, pro |
| DigitalOcean | In Progress | micro, small, pro |

## FinOps Dashboard

The web interface includes comprehensive FinOps capabilities:

- **Budget Tracking**: Real-time budget utilization by team and provider
- **Cost Estimation**: Accurate cost estimates for deployments before provisioning
- **Spend Analysis**: Historical cost data and trend analysis
- **Cost Alerts**: Automatic alerts when approaching budget thresholds
- **Provider Comparison**: Cost comparison across different cloud providers
- **Optimization Recommendations**: Suggestions for cost optimization

Access FinOps features through the web dashboard at `http://localhost:8080` or via CLI:
```bash
engine cost --team platform --month 2024-01
```

## Development

### Building from Source

```bash
# Install Go 1.21+
go version

# Clone repository
git clone https://github.com/nutcas3/the-engine.git
cd the-engine

# Install dependencies
make install-deps

# Build all binaries
make build-cli
make build-web
make encrypt

# Run tests
make test
```

### Project Structure
```
the-engine/
cmd/
  cli/           # CLI tool
  ui/            # Web server
  encrypt/       # Encryption tool
internal/
  handlers/      # HTTP handlers
  cache/         # Caching layer
  health/        # Health checks
  rate/          # Rate limiting
  config/        # Secure configuration
  kubernetes/    # Kubernetes client (internal use)
web/
  index.html     # Web UI with HTMX + Tailwind v4
  script.js      # Theme toggle and SSE
configs/         # Configuration files
```

## API Documentation

### Endpoints

- `GET /` - Web dashboard
- `GET /api/compositions` - List available compositions
- `GET /api/deployments` - List deployments
- `GET /api/cost/monthly?team=<team>` - Get cost data
- `GET /api/health/status` - Health check
- `GET /api/stream` - SSE real-time updates
- `GET /api/swagger` - OpenAPI documentation

### Rate Limiting
- 100 requests per second
- Burst: 10 requests
- Per-IP rate limiting

## Troubleshooting

### Installation Issues
```bash
# If download fails, check internet connection
curl -I https://github.com/nutcas3/the-engine/releases/

# If permissions error, use sudo
sudo ./install.sh
```

### Encryption Issues
```bash
# Ensure ENGINE_MASTER_KEY is set
echo $ENGINE_MASTER_KEY

# Regenerate if lost
openssl rand -base64 32
```

### Web Server Issues
```bash
# Check if port is in use
lsof -i :8080

# Change port in config.yaml
# web_port: 8081
```

## Contributing

1. Fork the repository
2. Create feature branch
3. Make changes
4. Add tests
5. Submit pull request

## License

MIT License - see [LICENSE](LICENSE) file for details.

---

*"Infrastructure is a liability. Control is a choice. The Engine is the bridge."*
