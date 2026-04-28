# Sovereign Engine: Crossplane Multi-Cloud Control Plane

A minimalist, high-performance Internal Developer Platform (IDP) built with Go, Crossplane, and HTMX. This architecture abstracts cloud complexity across AWS, Azure, GCP, Hetzner, OVH, and DigitalOcean while maintaining strict governance and cost control.

## Architecture

```
Developer Intent (XRD) 
    |
    v
Go Composition Function (The Brain)
    |
    v
[ FinOps Guardrails ] [ SLSA Security ] [ Policy Engine ]
    |
    v
Crossplane (Orchestration Layer)
    |
    v
[ AWS ] [ Azure ] [ GCP ] [ Hetzner ] [ OVH ] [ DigitalOcean ]
    |
    v
Cloudflare (Edge/DNS/WAF)
```

## Features

### Implemented Capabilities
- **Multi-Cloud Abstraction**: Single API for 6 cloud providers (AWS, Azure, GCP, Hetzner, OVH, DigitalOcean)
- **Real-time FinOps**: Budget validation and cost estimation
- **Provider Mapping**: Automatic tier-to-SKU translation for all clouds
- **CLI Interface**: Kamal-like command-line tool for deployments
- **Cost Guardrails**: Budget checking and automated optimization recommendations

### In Progress
- **Cloud Compositions**: AWS and Hetzner compositions complete (Azure, GCP, OVH, DigitalOcean pending)
- **UI Backend**: HTMX dashboard frontend complete (backend API in progress)
- **Security**: SLSA verification framework implemented (full integration pending)

### Planned Features
- **Self-Healing**: Automatic drift detection and reconciliation
- **Automated Thrift**: Intelligent tier downgrading for dev environments
- **Multi-Cloud Resilience**: Automatic failover between providers
- **Unified Observability**: LGTM stack integration
- **Authentication**: OAuth2/OIDC and API key management

## Quick Start

### Prerequisites
- Go 1.26+
- Docker (for container builds)
- k3s (lightweight Kubernetes) - optional for local development
- Crossplane - required for production deployment

### Local Development

```bash
# 1. Clone the repository
git clone https://github.com/nutcase/the-engine.git
cd the-engine

# 2. Install dependencies
make install-deps

# 3. Build the function binary
make build

# 4. Run the function locally (development mode)
make dev

# 5. Test the CLI
go run ./cmd/cli deploy --provider hetzner --tier micro --region nbg1
go run ./cmd/cli cost --team platform
go run ./cmd/cli drift-check
```

### Production Deployment

```bash
# 1. Install k3s (minimal footprint)
curl -sfL https://get.k3s.io | sh -s - --disable traefik --disable servicelb

# 2. Install Crossplane
helm repo add crossplane https://charts.crossplane.io/stable
helm install crossplane crossplane/crossplane --namespace crossplane-system --create-namespace

# 3. Install providers
kubectl apply -f configs/providers.yaml

# 4. Deploy the Engine API definitions
kubectl apply -f apis/

# 5. Deploy compositions (AWS and Hetzner currently available)
kubectl apply -f compositions/aws/
kubectl apply -f compositions/hetzner/

# 6. Build and deploy function container
docker build -f build/Dockerfile.function -t engine/function:latest .
```

## Usage

### CLI Interface

The CLI tool provides a Kamal-like interface for managing multi-cloud deployments.

```bash
# Deploy new instance (generates XRD manifest)
go run ./cmd/cli deploy --provider hetzner --tier micro --region nbg1

# Deploy with budget guardrails
go run ./cmd/cli deploy --provider azure --tier pro --budget 500 --team platform

# List all deployments (mock data for development)
go run ./cmd/cli list --provider hetzner

# Check cost report
go run ./cmd/cli cost --team platform

# Check for configuration drift (mock detection)
go run ./cmd/cli drift-check

# Get help
go run ./cmd/cli help
```

**Note**: The CLI currently generates XRD manifests and provides mock data for development. Full integration with Kubernetes and Crossplane is in progress.

### Web Dashboard

The HTMX dashboard frontend is implemented with the Vantablack/Sawdust sovereign theme. Backend API endpoints are currently being developed to provide real-time data for:
- Deployment status monitoring
- Budget utilization tracking
- Configuration drift alerts
- Quick deployment actions

## Architecture Components

### 1. Agnostic Contract (XRD)
Single API definition for all cloud providers:
```yaml
apiVersion: engine.io/v1alpha1
kind: XCompute
metadata:
  name: my-server
spec:
  provider: hetzner
  tier: micro
  region: nbg1
  budget_max: 100
```

### 2. Go Composition Function
The brain that maps intent to cloud resources with:
- Budget validation
- Security verification
- Provider mapping
- Cost optimization

### 3. Multi-Cloud Provider Mapping
| Tier | AWS | Azure | GCP | Hetzner | OVH | DigitalOcean |
|------|-----|-------|-----|---------|-----|-------------|
| micro | t3.micro | Standard_B1s | e2-micro | cx11 | s1-2 | s-1vcpu-1gb |
| small | t3.small | Standard_B2s | e2-small | cpx11 | s1-4 | s-1vcpu-2gb |
| pro | c6i.large | Standard_D2s_v5 | n2-standard-2 | cpx21 | s1-8 | s-2vcpu-4gb |

### 4. FinOps Guardrails
- **Budget Validation**: Check deployments against team budgets before provisioning
- **Cost Estimation**: Provide accurate cost estimates for different tiers and providers
- **Spend Tracking**: Monitor current spend across all cloud providers (mock data currently)
- **Cost Recommendations**: Automated suggestions for cost optimization

### 5. Security & Compliance
- **SLSA Framework**: Image verification framework implemented (Cosign integration pending)
- **Security Verification**: Security check functions implemented (full integration pending)
- **Policy Enforcement**: Policy validation framework ready (OPA integration planned)

## Development

### Project Structure
```
the-engine/
cmd/
  function/        # Crossplane Logic Engine
  cli/             # Kamal-like CLI tool
internal/
  provider/        # Multi-cloud mapping logic
  finops/          # Billing & Budget guardrails
  security/        # Cosign & Policy verification
apis/              # Agnostic XRD Definitions
ui/                # HTMX Dashboard
compositions/      # Crossplane compositions
build/             # Dockerfiles
configs/           # Provider configurations
```

### Build Commands
```bash
# Build optimized binary
make build

# Run tests
make test

# Clean build artifacts
make clean

# Deploy to cluster
make deploy
```

### Development Workflow
```bash
# Install dependencies
make install-deps

# Run function locally (development mode)
make dev

# Build optimized function binary
make build

# Clean build artifacts
make clean

# Deploy to Kubernetes (when compositions are ready)
make deploy
```

### Testing
```bash
# Run tests (test suite being developed)
make test

# Test FinOps functions
go test ./internal/finops/

# Test provider mapping
go test ./internal/provider/

# Test CLI commands
go run ./cmd/cli deploy --provider hetzner --tier micro --region nbg1
```

## Configuration

### Environment Variables
```bash
# Cloud Provider Credentials
export AWS_ACCESS_KEY_ID=...
export AWS_SECRET_ACCESS_KEY=...
export AZURE_CLIENT_ID=...
export AZURE_CLIENT_SECRET=...
export HCLOUD_TOKEN=...

# Feature Flags
export ENABLE_FINOPS=true
export ENABLE_SECURITY=true
export LOG_LEVEL=info
```

### Team Budgets
```yaml
# budgets.yaml
teams:
  platform:
    monthly_budget: 2000
    alert_thresholds: [50, 75, 90]
  engineering:
    monthly_budget: 1500
    alert_thresholds: [60, 80, 95]
```

## Monitoring & Observability

### Planned LGTM Stack Integration
The platform is designed to integrate with the LGTM stack for comprehensive observability:
- **Loki**: Centralized logs from all clouds
- **Grafana**: Unified dashboards for multi-cloud metrics
- **Tempo**: Distributed tracing across providers
- **Mimir**: Scalable metrics storage

**Status**: Framework ready, integration planned for Phase 2

### Key Metrics (Planned)
- `engine_deployments_total` - Deployment count by provider
- `engine_cost_monthly` - Monthly spend by team
- `engine_drift_events` - Configuration drift occurrences
- `engine_security_violations` - SLSA policy violations

## Security

### Supply Chain Security (Planned)
- **Cosign**: Image signing and verification (framework implemented, integration pending)
- **SBOM**: Software Bill of Materials generation (planned for Phase 5)
- **Attestations**: Cryptographic proof of build process (planned)

### Network Security (Planned)
- **Cloudflare WAF**: Edge protection for all endpoints (configuration planned)
- **Zero Trust**: JIT access with 1-hour expiration (planned for Phase 3)
- **Private Networking**: VPC peering across providers (planned)

## Troubleshooting

### Common Issues

**Build Errors**
```bash
# If go build fails, try cleaning and rebuilding
make clean
make install-deps
make build
```

**CLI Not Working**
```bash
# Ensure you're running from the project root
cd /path/to/the-engine
go run ./cmd/cli help

# Check Go version (requires 1.26+)
go version
```

**Crossplane Provider Issues**
```bash
# Check provider status
kubectl get providers

# View provider logs
kubectl logs -n crossplane-system deployment/provider-aws-ec2
```

**Function Deployment Issues**
```bash
# Check if function is running
kubectl get pods -n engine-system

# View function logs
kubectl logs -n engine-system deployment/engine-function

# Test with dry-run
kubectl apply --dry-run=client -f apis/xcompute.yaml
```

## Performance

### Optimization Techniques (Implemented)
- **Binary Size**: `-ldflags="-s -w"` and UPX compression for minimal footprint
- **Memory**: k3s with minimal components (~100MB footprint)
- **Efficient Mapping**: Direct tier-to-SKU translation without complex logic

### Planned Optimizations
- **Network**: Local registry for remote regions
- **Caching**: Resource caching in Go function
- **Connection Pooling**: Database and API connection optimization

### Benchmarks (Target)
| Operation | Target Latency (p50) | Target Latency (p99) |
|-----------|---------------------|---------------------|
| Deploy XCompute | 2.3s | 5.1s |
| Drift Detection | 500ms | 1.2s |
| Cost Query | 200ms | 800ms |
| Health Check | 50ms | 150ms |

**Status**: Performance optimization planned for Phase 4

## Contributing

We welcome contributions! Here's how to get started:

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Run tests (`make test`)
4. Build (`make build`)
5. Commit changes (`git commit -m 'Add amazing feature'`)
6. Push to branch (`git push origin feature/amazing-feature`)
7. Open Pull Request

### Contribution Guidelines
- Follow the existing code style and patterns
- Add tests for new features
- Update documentation as needed
- Check the [TODOs.md](TODOs.md) for prioritized tasks
- Focus on high-priority items first (missing cloud compositions, UI backend, etc.)

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Support

- **Issues**: https://github.com/your-org/the-engine/issues
- **Discussions**: https://github.com/your-org/the-engine/discussions
- **Roadmap**: See [TODOs.md](TODOs.md) for detailed implementation plan
- **Documentation**: See [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md) for project overview

## Current Status

**Completed:**
- ✅ Multi-cloud provider mapping (AWS, Azure, GCP, Hetzner, OVH, DigitalOcean)
- ✅ FinOps budget validation and cost estimation
- ✅ CLI tool with deploy, list, cost, and drift-check commands
- ✅ Go composition function with budget guardrails
- ✅ XRD API definition
- ✅ AWS and Hetzner Crossplane compositions
- ✅ Security verification framework
- ✅ Optimized build pipeline with UPX compression

**In Progress:**
- 🔄 UI backend API endpoints
- 🔄 Additional cloud compositions (Azure, GCP, OVH, DigitalOcean)
- 🔄 Kubernetes deployment manifests
- 🔄 Testing framework

**Planned:**
- 📋 Authentication and authorization
- 📋 Database persistence
- 📋 LGTM stack integration
- 📋 CI/CD pipeline
- 📋 Advanced security features

---

*"Infrastructure is a liability. Control is a choice. The Engine is the bridge."*
