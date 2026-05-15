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
