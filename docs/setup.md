# Setup Guide

This guide covers installation, configuration, and initial setup of Sovereign Engine.

## Quick Start Installation

### Prerequisites

- curl or wget (for downloading binaries)
- No Go, Kubernetes, or Docker required for CLI usage

### Automated Installation

```bash
# Download and install pre-built binaries
curl -fsSL https://raw.githubusercontent.com/nutcas3/the-engine/main/install.sh | bash

# Set your master key (generated during install)
export ENGINE_MASTER_KEY=<your_generated_key>
```

### Manual Installation

#### Linux

```bash
wget https://github.com/nutcas3/the-engine/releases/download/v1.0.0/engine-linux-amd64.tar.gz
tar -xzf engine-linux-amd64.tar.gz
sudo cp engine /usr/local/bin/
sudo cp engine-encrypt /usr/local/bin/
sudo cp engine-web /usr/local/bin/
```

#### macOS

```bash
wget https://github.com/nutcas3/the-engine/releases/download/v1.0.0/engine-darwin-amd64.tar.gz
tar -xzf engine-darwin-amd64.tar.gz
sudo cp engine /usr/local/bin/
sudo cp engine-encrypt /usr/local/bin/
sudo cp engine-web /usr/local/bin/
```

## Configuration

### Initial Configuration

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

Generate and configure encryption:

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

## Cloud Provider Setup

### AWS

```bash
# Encrypt AWS credentials
engine-encrypt -key AWS_ACCESS_KEY_ID -value 'AKIA...' -file ~/.engine/secure.env
engine-encrypt -key AWS_SECRET_ACCESS_KEY -value 'secret...' -file ~/.engine/secure.env
engine-encrypt -key AWS_REGION -value 'us-east-1' -file ~/.engine/secure.env

# Test connection
engine status
```

### Azure

```bash
# Encrypt Azure credentials
engine-encrypt -key AZURE_CLIENT_ID -value 'client-id' -file ~/.engine/secure.env
engine-encrypt -key AZURE_CLIENT_SECRET -value 'secret' -file ~/.engine/secure.env
engine-encrypt -key AZURE_TENANT_ID -value 'tenant-id' -file ~/.engine/secure.env
engine-encrypt -key AZURE_SUBSCRIPTION_ID -value 'sub-id' -file ~/.engine/secure.env
```

### GCP

```bash
# Encrypt GCP credentials
engine-encrypt -key GCP_PROJECT_ID -value 'project-id' -file ~/.engine/secure.env
engine-encrypt -key GCP_CREDENTIALS -value 'path/to/credentials.json' -file ~/.engine/secure.env
```

### Hetzner

```bash
# Encrypt Hetzner credentials
engine-encrypt -key HETZNER_API_TOKEN -value 'token' -file ~/.engine/secure.env
```

### OVH

```bash
# Encrypt OVH credentials
engine-encrypt -key OVH_APPLICATION_KEY -value 'key' -file ~/.engine/secure.env
engine-encrypt -key OVH_APPLICATION_SECRET -value 'secret' -file ~/.engine/secure.env
engine-encrypt -key OVH_CONSUMER_KEY -value 'consumer-key' -file ~/.engine/secure.env
engine-encrypt -key OVH_ENDPOINT -value 'ovh-eu' -file ~/.engine/secure.env
```

### DigitalOcean

```bash
# Encrypt DigitalOcean credentials
engine-encrypt -key DIGITALOCEAN_TOKEN -value 'token' -file ~/.engine/secure.env
```

## Verification

### Test Installation

```bash
# Check CLI version
engine --version

# Check health status
engine status

# List available providers
engine list --providers
```

### Test Web Interface

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

## Production Setup

### Environment Variables

For production deployments, configure the following environment variables:

```bash
# Master key (required)
export ENGINE_MASTER_KEY=<your_master_key>

# Web server port (default: 8080)
export ENGINE_WEB_PORT=8080

# Function runner port (default: 9443)
export FUNCTION_LISTEN_ADDRESS=:9443
```

### Production Integrations

Configure monitoring, CI/CD, and deployment tracking:

```bash
# Monitoring system (prometheus, datadog, cloudwatch)
export MONITORING_SYSTEM=prometheus
export PROMETHEUS_URL=https://prometheus.example.com
export PROMETHEUS_TOKEN=your_token

# CI/CD system (github, jenkins, gitlab)
export CICD_SYSTEM=github
export GITHUB_TOKEN=ghp_xxxxx
export GITHUB_REPOSITORY=owner/repo

# Deployment system (argocd, flux, kubernetes)
export DEPLOYMENT_SYSTEM=argocd
export ARGOCD_URL=https://argocd.example.com
export ARGOCD_TOKEN=your_token
```

See the [Monitoring Guide](monitoring.md), [CI/CD Guide](cicd.md), and [Deployment Guide](deployments.md) for detailed configuration.

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

### Permission Issues

```bash
# Ensure secure env file has correct permissions
chmod 600 ~/.engine/secure.env

# Ensure config directory has correct permissions
chmod 700 ~/.engine
```

## Next Steps

- Configure [Monitoring Integration](monitoring.md)
- Set up [CI/CD Integration](cicd.md)
- Configure [Deployment Tracking](deployments.md)
- Define [Cleanup Policies](cleanup.md)
