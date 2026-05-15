# CI/CD Integration with Sovereign Engine

This guide covers integrating Sovereign Engine into CI/CD pipelines for automated environment management, cleanup, and deployment tracking.

## Overview

Sovereign Engine can be integrated into CI/CD workflows to:
- Automatically clean up test environments after tests complete
- Trigger environment shutdown before deployment
- Monitor deployment status and health
- Manage resource lifecycle across pipelines

## GitHub Actions Integration

### Basic Setup

Add Sovereign Engine to your GitHub Actions workflow:

```yaml
name: CI/CD Pipeline

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]
  workflow_dispatch:

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Install Sovereign Engine
        run: |
          curl -fsSL https://raw.githubusercontent.com/nutcas3/the-engine/main/install.sh | bash
          echo "$HOME/.engine" >> $GITHUB_PATH

      - name: Configure Engine
        env:
          ENGINE_MASTER_KEY: ${{ secrets.ENGINE_MASTER_KEY }}
        run: |
          export ENGINE_MASTER_KEY=$ENGINE_MASTER_KEY
          engine status

      - name: Deploy Test Environment
        run: |
          engine deploy --environment test --provider aws

      - name: Run Tests
        run: |
          npm test

      - name: Cleanup Test Environment
        if: always()
        run: |
          engine cleanup --environment test --force
```

### Test Environment Management

Automated test environment lifecycle:

```yaml
name: Test Pipeline

on:
  push:
    branches: [ develop ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Install Sovereign Engine
        run: |
          curl -fsSL https://raw.githubusercontent.com/nutcas3/the-engine/main/install.sh | bash
          echo "$HOME/.engine" >> $GITHUB_PATH

      - name: Configure Engine
        env:
          ENGINE_MASTER_KEY: ${{ secrets.ENGINE_MASTER_KEY }}
          AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
          AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
        run: |
          export ENGINE_MASTER_KEY=$ENGINE_MASTER_KEY
          engine-encrypt -key AWS_ACCESS_KEY_ID -value "$AWS_ACCESS_KEY_ID" -file ~/.engine/secure.env
          engine-encrypt -key AWS_SECRET_ACCESS_KEY -value "$AWS_SECRET_ACCESS_KEY" -file ~/.engine/secure.env

      - name: Create Test Environment
        run: |
          engine deploy --environment test --provider aws --composition test-app

      - name: Run Integration Tests
        run: |
          npm run test:integration

      - name: Check Test Status
        run: |
          engine status --environment test

      - name: Cleanup on Success
        if: success()
        run: |
          engine cleanup --environment test --auto-nuke

      - name: Keep Environment on Failure
        if: failure()
        run: |
          engine notify --environment test --message "Tests failed, environment preserved for debugging"
```

### Staging Deployment with Cleanup

Deploy to staging and clean up old environments:

```yaml
name: Deploy to Staging

on:
  push:
    branches: [ main ]

jobs:
  deploy-staging:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Install Sovereign Engine
        run: |
          curl -fsSL https://raw.githubusercontent.com/nutcas3/the-engine/main/install.sh | bash
          echo "$HOME/.engine" >> $GITHUB_PATH

      - name: Configure Engine
        env:
          ENGINE_MASTER_KEY: ${{ secrets.ENGINE_MASTER_KEY }}
        run: |
          export ENGINE_MASTER_KEY=$ENGINE_MASTER_KEY
          engine status

      - name: Cleanup Old Staging Environments
        run: |
          engine cleanup --environment staging --nuke-old --keep-latest 2

      - name: Deploy to Staging
        run: |
          engine deploy --environment staging --provider aws --composition staging-app

      - name: Verify Deployment
        run: |
          engine verify --environment staging --health-check

      - name: Run Smoke Tests
        run: |
          npm run test:smoke
```

### Multi-Environment Pipeline

Manage multiple environments in a single pipeline:

```yaml
name: Multi-Environment Pipeline

on:
  push:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Install Sovereign Engine
        run: |
          curl -fsSL https://raw.githubusercontent.com/nutcas3/the-engine/main/install.sh | bash
          echo "$HOME/.engine" >> $GITHUB_PATH

      - name: Configure Engine
        env:
          ENGINE_MASTER_KEY: ${{ secrets.ENGINE_MASTER_KEY }}
        run: |
          export ENGINE_MASTER_KEY=$ENGINE_MASTER_KEY
          engine status

      - name: Deploy Test Environment
        run: |
          engine deploy --environment test --provider aws

      - name: Run Tests
        run: |
          npm test

      - name: Cleanup Test Environment
        if: always()
        run: |
          engine cleanup --environment test --force

  deploy-staging:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Install Sovereign Engine
        run: |
          curl -fsSL https://raw.githubusercontent.com/nutcas3/the-engine/main/install.sh | bash
          echo "$HOME/.engine" >> $GITHUB_PATH

      - name: Deploy to Staging
        run: |
          engine deploy --environment staging --provider aws

      - name: Verify Staging
        run: |
          engine verify --environment staging

  deploy-prod:
    needs: deploy-staging
    runs-on: ubuntu-latest
    environment: production
    steps:
      - uses: actions/checkout@v3

      - name: Install Sovereign Engine
        run: |
          curl -fsSL https://raw.githubusercontent.com/nutcas3/the-engine/main/install.sh | bash
          echo "$HOME/.engine" >> $GITHUB_PATH

      - name: Deploy to Production
        run: |
          engine deploy --environment prod --provider aws

      - name: Verify Production
        run: |
          engine verify --environment prod --health-check
```

## Jenkins Integration

### Basic Jenkinsfile

```groovy
pipeline {
    agent any

    environment {
        ENGINE_MASTER_KEY = credentials('ENGINE_MASTER_KEY')
        AWS_ACCESS_KEY_ID = credentials('AWS_ACCESS_KEY_ID')
        AWS_SECRET_ACCESS_KEY = credentials('AWS_SECRET_ACCESS_KEY')
    }

    stages {
        stage('Install Engine') {
            steps {
                sh '''
                    curl -fsSL https://raw.githubusercontent.com/nutcas3/the-engine/main/install.sh | bash
                    export PATH="$HOME/.engine:$PATH"
                '''
            }
        }

        stage('Configure Engine') {
            steps {
                sh '''
                    export ENGINE_MASTER_KEY=${ENGINE_MASTER_KEY}
                    engine-encrypt -key AWS_ACCESS_KEY_ID -value "${AWS_ACCESS_KEY_ID}" -file ~/.engine/secure.env
                    engine-encrypt -key AWS_SECRET_ACCESS_KEY -value "${AWS_SECRET_ACCESS_KEY}" -file ~/.engine/secure.env
                    engine status
                '''
            }
        }

        stage('Deploy Test') {
            steps {
                sh 'engine deploy --environment test --provider aws'
            }
        }

        stage('Run Tests') {
            steps {
                sh 'npm test'
            }
        }

        stage('Cleanup Test') {
            steps {
                sh 'engine cleanup --environment test --force'
            }
        }
    }

    post {
        always {
            sh 'engine status --environment test'
        }
    }
}
```

### Multi-Stage Jenkins Pipeline

```groovy
pipeline {
    agent any

    environment {
        ENGINE_MASTER_KEY = credentials('ENGINE_MASTER_KEY')
    }

    stages {
        stage('Install Engine') {
            steps {
                sh 'curl -fsSL https://raw.githubusercontent.com/nutcas3/the-engine/main/install.sh | bash'
                sh 'export PATH="$HOME/.engine:$PATH"'
            }
        }

        stage('Test') {
            stages {
                stage('Deploy Test') {
                    steps {
                        sh 'engine deploy --environment test --provider aws'
                    }
                }

                stage('Run Tests') {
                    steps {
                        sh 'npm test'
                    }
                }

                stage('Cleanup Test') {
                    steps {
                        sh 'engine cleanup --environment test --force'
                    }
                }
            }
        }

        stage('Staging') {
            stages {
                stage('Deploy Staging') {
                    steps {
                        sh 'engine deploy --environment staging --provider aws'
                    }
                }

                stage('Verify Staging') {
                    steps {
                        sh 'engine verify --environment staging'
                    }
                }
            }
        }

        stage('Production') {
            when {
                branch 'main'
            }
            stages {
                stage('Deploy Production') {
                    steps {
                        sh 'engine deploy --environment prod --provider aws'
                    }
                }

                stage('Verify Production') {
                    steps {
                        sh 'engine verify --environment prod --health-check'
                    }
                }
            }
        }
    }

    post {
        failure {
            sh 'engine notify --environment test --message "Pipeline failed"'
        }
    }
}
```

## GitLab CI Integration

### Basic GitLab CI

```yaml
stages:
  - test
  - deploy
  - cleanup

variables:
  ENGINE_MASTER_KEY: $ENGINE_MASTER_KEY

before_script:
  - curl -fsSL https://raw.githubusercontent.com/nutcas3/the-engine/main/install.sh | bash
  - export PATH="$HOME/.engine:$PATH"
  - export ENGINE_MASTER_KEY=$ENGINE_MASTER_KEY

test:
  stage: test
  script:
    - engine deploy --environment test --provider aws
    - npm test
    - engine cleanup --environment test --force
  only:
    - merge_requests
    - develop

deploy-staging:
  stage: deploy
  script:
    - engine deploy --environment staging --provider aws
    - engine verify --environment staging
  only:
    - main

cleanup-staging:
  stage: cleanup
  script:
    - engine cleanup --environment staging --nuke-old --keep-latest 2
  only:
    - main
  when: manual
```

### Advanced GitLab CI Pipeline

```yaml
stages:
  - install
  - test
  - staging
  - production
  - cleanup

variables:
  ENGINE_MASTER_KEY: $ENGINE_MASTER_KEY
  AWS_ACCESS_KEY_ID: $AWS_ACCESS_KEY_ID
  AWS_SECRET_ACCESS_KEY: $AWS_SECRET_ACCESS_KEY

install-engine:
  stage: install
  script:
    - curl -fsSL https://raw.githubusercontent.com/nutcas3/the-engine/main/install.sh | bash
    - export PATH="$HOME/.engine:$PATH"
    - export ENGINE_MASTER_KEY=$ENGINE_MASTER_KEY
    - engine-encrypt -key AWS_ACCESS_KEY_ID -value "$AWS_ACCESS_KEY_ID" -file ~/.engine/secure.env
    - engine-encrypt -key AWS_SECRET_ACCESS_KEY -value "$AWS_SECRET_ACCESS_KEY" -file ~/.engine/secure.env
    - engine status

test:
  stage: test
  dependencies:
    - install-engine
  script:
    - engine deploy --environment test --provider aws
    - npm test
    - engine cleanup --environment test --force
  only:
    - merge_requests
    - develop

deploy-staging:
  stage: staging
  dependencies:
    - install-engine
  script:
    - engine cleanup --environment staging --nuke-old --keep-latest 1
    - engine deploy --environment staging --provider aws
    - engine verify --environment staging
  only:
    - main

deploy-production:
  stage: production
  dependencies:
    - install-engine
  script:
    - engine deploy --environment prod --provider aws
    - engine verify --environment prod --health-check
  only:
    - main
  when: manual

cleanup-all:
  stage: cleanup
  dependencies:
    - install-engine
  script:
    - engine cleanup --environment test --force
    - engine cleanup --environment staging --nuke-old --keep-latest 1
  only:
    - schedules
```

## Common CI/CD Patterns

### Test Environment Lifecycle

```yaml
# Create, test, cleanup pattern
steps:
  - name: Deploy Test Environment
    run: engine deploy --environment test --provider aws

  - name: Run Tests
    run: npm test

  - name: Cleanup on Success
    if: success()
    run: engine cleanup --environment test --auto-nuke

  - name: Keep on Failure
    if: failure()
    run: engine notify --environment test --message "Tests failed"
```

### Blue-Green Deployment

```yaml
# Blue-green deployment with rollback
steps:
  - name: Deploy Blue
    run: engine deploy --environment staging-blue --provider aws

  - name: Verify Blue
    run: engine verify --environment staging-blue

  - name: Switch Traffic to Blue
    run: engine switch --environment staging --target blue

  - name: Deploy Green
    run: engine deploy --environment staging-green --provider aws

  - name: Verify Green
    run: engine verify --environment staging-green

  - name: Switch Traffic to Green
    run: engine switch --environment staging --target green

  - name: Cleanup Blue
    run: engine cleanup --environment staging-blue --force
```

### Canary Deployment

```yaml
# Canary deployment with monitoring
steps:
  - name: Deploy Canary
    run: engine deploy --environment staging-canary --provider aws --scale 10%

  - name: Monitor Canary
    run: |
      for i in {1..30}; do
        engine monitor --environment staging-canary
        sleep 60
      done

  - name: Full Deployment
    run: engine deploy --environment staging --provider aws

  - name: Cleanup Canary
    run: engine cleanup --environment staging-canary --force
```

### Rollback on Failure

```yaml
# Automatic rollback on verification failure
steps:
  - name: Deploy
    run: engine deploy --environment staging --provider aws

  - name: Verify
    id: verify
    run: |
      if ! engine verify --environment staging; then
        echo "Verification failed, rolling back"
        exit 1
      fi

  - name: Rollback on Failure
    if: failure()
    run: engine rollback --environment staging
```

## Secrets Management

### GitHub Actions Secrets

Configure these secrets in your repository settings:

```
ENGINE_MASTER_KEY
AWS_ACCESS_KEY_ID
AWS_SECRET_ACCESS_KEY
GITHUB_TOKEN
PROMETHEUS_TOKEN
DATADOG_API_KEY
```

### Jenkins Credentials

Add credentials in Jenkins:

```
ENGINE_MASTER_KEY - Secret text
AWS_ACCESS_KEY_ID - Secret text
AWS_SECRET_ACCESS_KEY - Secret text
```

### GitLab CI Variables

Add variables in GitLab CI/CD settings:

```
ENGINE_MASTER_KEY - Masked, Protected
AWS_ACCESS_KEY_ID - Masked, Protected
AWS_SECRET_ACCESS_KEY - Masked, Protected
```

## CLI Commands in CI/CD

### Deployment Commands

```bash
# Deploy environment
engine deploy --environment <env> --provider <provider> --composition <name>

# Verify deployment
engine verify --environment <env> --health-check

# Rollback deployment
engine rollback --environment <env>

# Switch traffic (blue-green)
engine switch --environment <env> --target <blue|green>
```

### Cleanup Commands

```bash
# Cleanup environment
engine cleanup --environment <env> --force

# Auto-nuke after threshold
engine cleanup --environment <env> --auto-nuke

# Nuke old environments
engine cleanup --environment <env> --nuke-old --keep-latest <n>

# Manual nuke
engine nuke --environment <env> --provider <provider>
```

### Status Commands

```bash
# Check environment status
engine status --environment <env>

# List resources
engine list --provider <provider> --environment <env>

# Monitor environment
engine monitor --environment <env>
```

### Notification Commands

```bash
# Send notification
engine notify --environment <env> --message "Message"

# Alert on failure
engine alert --environment <env> --threshold <value>
```

## Best Practices

1. **Always cleanup on success** - Clean up test environments after successful tests
2. **Keep environments on failure** - Preserve failed environments for debugging
3. **Use environment variables** - Store secrets in CI/CD secrets, not in code
4. **Verify deployments** - Always verify deployment health before proceeding
5. **Monitor resource usage** - Track costs and resource consumption
6. **Use cleanup policies** - Configure automatic cleanup policies
7. **Test locally first** - Test engine commands locally before CI/CD integration
8. **Use manual approvals** - Require manual approval for production deployments
9. **Implement rollbacks** - Always have a rollback strategy
10. **Log everything** - Enable logging for debugging and auditing

## Troubleshooting

### Engine Not Found

```bash
# Ensure engine is installed and in PATH
export PATH="$HOME/.engine:$PATH"
which engine
```

### Authentication Issues

```bash
# Verify master key is set
echo $ENGINE_MASTER_KEY

# Test engine status
engine status
```

### Permission Issues

```bash
# Verify credentials are configured
engine-encrypt -key TEST -value "test" -file ~/.engine/secure.env
```

### Cleanup Not Working

```bash
# Check cleanup policies
engine policies --environment <env>

# Manually trigger cleanup
engine cleanup --environment <env> --force
```

## Next Steps

- Configure [Monitoring Integration](monitoring.md)
- Set up [Deployment Tracking](deployments.md)
- Define [Cleanup Policies](cleanup.md)
- Review [CI/CD Integration](cicd.md)
