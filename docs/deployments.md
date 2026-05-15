# Deployment Tracking Guide

This guide covers configuring deployment tracking systems for detecting recent deployments before cleanup operations.

## Overview

The cleanup system integrates with deployment platforms to:
- Detect recent deployments before nuking dev environments
- Track deployment timestamps and sync status
- Integrate with ArgoCD, Flux, and Kubernetes

## Supported Deployment Systems

- **ArgoCD** - GitOps deployment system with application sync tracking
- **Flux** - GitOps toolkit with Kustomization and HelmRelease reconciliation
- **Kubernetes** - Native Kubernetes deployment tracking

## Configuration

### ArgoCD

#### Environment Variables

```bash
export DEPLOYMENT_SYSTEM=argocd
export ARGOCD_URL=https://argocd.example.com
export ARGOCD_TOKEN=your_token
```

#### Token Setup

1. Generate an ArgoCD API Token:
   ```bash
   argocd account generate-token --account <username>
   ```

2. Set the environment variable:
   ```bash
   export ARGOCD_TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
   ```

#### Application Labeling

ArgoCD applications must be labeled with the environment:
```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: myapp-dev
  labels:
    environment: dev
spec:
  # ... application spec
```

#### Example Setup

```bash
# Configure ArgoCD
export DEPLOYMENT_SYSTEM=argocd
export ARGOCD_URL=https://argocd.example.com
export ARGOCD_TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

# Test configuration
engine status
```

#### API Endpoints

The system uses the ArgoCD REST API:
- Get applications: `GET /api/v1/applications?selector=environment={env}`
- Check sync status and timestamps

### Flux

#### Environment Variables

```bash
export DEPLOYMENT_SYSTEM=flux
```

#### Kubernetes Configuration

Flux uses the Kubernetes client for accessing Flux CRDs:

```bash
# In-cluster configuration (default)
# Uses service account when running in Kubernetes

# Local development
export KUBECONFIG=/path/to/kubeconfig
```

#### Resource Labeling

Flux resources must be labeled with the environment:
```yaml
apiVersion: kustomize.toolkit.fluxcd.io/v1
kind: Kustomization
metadata:
  name: myapp-dev
  labels:
    environment: dev
spec:
  # ... kustomization spec
```

```yaml
apiVersion: helm.toolkit.fluxcd.io/v2
kind: HelmRelease
metadata:
  name: myapp-dev
  labels:
    environment: dev
spec:
  # ... helm release spec
```

#### Example Setup

```bash
# Configure Flux
export DEPLOYMENT_SYSTEM=flux

# Configure Kubernetes access (if not in-cluster)
export KUBECONFIG=/path/to/kubeconfig

# Test configuration
engine status
```

#### CRDs Accessed

The system queries Flux CRDs:
- `kustomize.toolkit.fluxcd.io/v1/kustomizations`
- `helm.toolkit.fluxcd.io/v2/helmreleases`

### Kubernetes

#### Environment Variables

```bash
export DEPLOYMENT_SYSTEM=kubernetes
```

#### Kubernetes Configuration

Kubernetes deployment tracking uses the Kubernetes client:

```bash
# In-cluster configuration (default)
# Uses service account when running in Kubernetes

# Local development
export KUBECONFIG=/path/to/kubeconfig
```

#### Deployment Labeling

Kubernetes deployments must be labeled with the environment:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp-dev
  labels:
    environment: dev
spec:
  # ... deployment spec
```

#### Example Setup

```bash
# Configure Kubernetes
export DEPLOYMENT_SYSTEM=kubernetes

# Configure Kubernetes access (if not in-cluster)
export KUBECONFIG=/path/to/kubeconfig

# Test configuration
engine status
```

#### Resources Accessed

The system queries Kubernetes resources:
- `apps/v1/deployments` across all namespaces

## How It Works

### Recent Deployment Check

The `checkNoRecentDeployments` function determines if there have been recent deployments:

```go
noRecent := cm.checkNoRecentDeployments(ctx, "dev", 24*time.Hour)
```

This checks for:
- **ArgoCD**: Application sync timestamps within threshold
- **Flux**: Kustomization/HelmRelease reconciliation timestamps
- **Kubernetes**: Deployment creation/update timestamps

### ArgoCD Sync Tracking

For ArgoCD, the system checks:
- Application sync status
- Last sync finished timestamp
- Environment-labeled applications

### Flux Reconciliation Tracking

For Flux, the system checks:
- Kustomization last applied revision
- Kustomization last attempted revision
- HelmRelease last applied revision
- Condition transition timestamps

### Kubernetes Deployment Tracking

For Kubernetes, the system checks:
- Deployment creation timestamp
- Deployment update timestamps
- Condition last transition times

## Integration with Cleanup Policies

### Dev Environments

Dev environments check for recent deployments before nuking:

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
1. **No deployments in last 24 hours** (via deployment tracking)
2. No active resources
3. No user activity (via monitoring)

### Test Environments

Test environments can also use deployment tracking:

```yaml
test-auto-cleanup:
  environment: test
  auto_shutdown: true
  shutdown_after: 2h
  nuke_after: 6h
  exclude_patterns: []
```

## Troubleshooting

### ArgoCD Issues

```bash
# Test ArgoCD connectivity
curl -H "Authorization: Bearer $ARGOCD_TOKEN" \
  https://argocd.example.com/api/v1/applications

# Check for environment-labeled applications
curl -H "Authorization: Bearer $ARGOCD_TOKEN" \
  "https://argocd.example.com/api/v1/applications?selector=environment=dev"
```

### Flux Issues

```bash
# Check Kubernetes connectivity
kubectl get kustomizations -A
kubectl get helmreleases -A

# Check for environment-labeled resources
kubectl get kustomizations -l environment=dev -A
kubectl get helmreleases -l environment=dev -A
```

### Kubernetes Issues

```bash
# Check Kubernetes connectivity
kubectl get deployments -A

# Check for environment-labeled deployments
kubectl get deployments -l environment=dev -A
```

### Deployment Tracking Not Working

If deployment checks always return false (deployments considered recent):

1. Verify environment variables are set
2. Test API connectivity (ArgoCD) or Kubernetes access
3. Check for proper resource labeling
4. Verify permissions/credentials
5. Check for recent deployments

## Best Practices

1. **Label all resources** - Ensure all deployments/CRDs have environment labels
2. **Use consistent naming** - Use consistent naming conventions across environments
3. **Monitor deployment health** - Ensure your deployment system is reliable
4. **Test with staging** - Test cleanup policies in staging first
5. **Use webhook notifications** - Get notified of deployment failures

## Security Considerations

- Use service accounts with minimal permissions
- Rotate API tokens regularly
- Use read-only tokens where possible
- Store secrets securely (use `engine-encrypt`)
- Monitor for unauthorized API access
- Use RBAC in Kubernetes for least privilege

## Examples

### ArgoCD with Multi-Environment Applications

```bash
# Configure
export DEPLOYMENT_SYSTEM=argocd
export ARGOCD_URL=https://argocd.example.com
export ARGOCD_TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

# The system will:
# 1. Query applications with environment label
# 2. Check sync status and timestamps
# 3. Prevent nuke if sync occurred recently
```

### Flux with Kustomizations

```bash
# Configure
export DEPLOYMENT_SYSTEM=flux
export KUBECONFIG=/path/to/kubeconfig

# The system will:
# 1. Query Kustomizations with environment label
# 2. Check reconciliation timestamps
# 3. Prevent nuke if reconciliation occurred recently
```

### Kubernetes with Native Deployments

```bash
# Configure
export DEPLOYMENT_SYSTEM=kubernetes
export KUBECONFIG=/path/to/kubeconfig

# The system will:
# 1. Query deployments with environment label
# 2. Check creation/update timestamps
# 3. Prevent nuke if deployment occurred recently
```

## Kubernetes Client Configuration

### In-Cluster Configuration

When running in Kubernetes, the system uses the service account:

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: engine-cleanup
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: engine-cleanup
rules:
- apiGroups: ["kustomize.toolkit.fluxcd.io"]
  resources: ["kustomizations"]
  verbs: ["get", "list"]
- apiGroups: ["helm.toolkit.fluxcd.io"]
  resources: ["helmreleases"]
  verbs: ["get", "list"]
- apiGroups: ["apps"]
  resources: ["deployments"]
  verbs: ["get", "list"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: engine-cleanup
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: engine-cleanup
subjects:
- kind: ServiceAccount
  name: engine-cleanup
  namespace: default
```

### Local Development

For local development, use kubeconfig:

```bash
export KUBECONFIG=/path/to/kubeconfig
```

## Next Steps

- Configure [Monitoring Integration](monitoring.md)
- Set up [CI/CD Integration](cicd.md)
- Define [Cleanup Policies](cleanup.md)
