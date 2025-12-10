# 2 Environment Setup: Develop (Staging) & Production

## Overview

This project uses **2 environments** for deployment:

1. **Staging/Development** - Branch: `develop`
2. **Production** - Branch: `main`

## How It Works

### GitHub Branches → Environments

```
┌─────────────┐         ┌──────────────────┐         ┌─────────────┐
│  develop    │ ──────> │  GitHub Actions  │ ──────> │   Staging   │
│   branch    │         │   CI/CD Deploy   │         │ Environment │
└─────────────┘         └──────────────────┘         └─────────────┘

┌─────────────┐         ┌──────────────────┐         ┌─────────────┐
│    main     │ ──────> │  GitHub Actions  │ ──────> │ Production  │
│   branch    │         │   CI/CD Deploy   │         │ Environment │
└─────────────┘         └──────────────────┘         └─────────────┘
```

### Workflow Triggers

| Branch | Environment | URL | Cluster | Replicas |
|--------|-------------|-----|---------|----------|
| `develop` | Staging | https://staging.suitemedia.com | EKS-staging | 2-5 pods |
| `main` | Production | https://suitemedia.com | EKS-production | 6-20 pods |

## AWS Secrets Manager Structure

### Staging Environment
**Secret Name:** `staging/suitemedia/app`

```json
{
  "DB_HOST": "staging-aurora.cluster-xxx.rds.amazonaws.com",
  "DB_PORT": "5432",
  "DB_USER": "suitemedia_user",
  "DB_PASSWORD": "staging-password-here",
  "DB_NAME": "suitemedia_staging",
  "DB_SSLMODE": "require",
  "REDIS_HOST": "staging-redis.cache.amazonaws.com",
  "REDIS_PORT": "6379",
  "REDIS_PASSWORD": "staging-redis-pass",
  "JWT_SECRET": "staging-jwt-secret-key-change-me",
  "JWT_REFRESH_SECRET": "staging-refresh-secret-key",
  "JWT_EXPIRATION_HOURS": "24",
  "JWT_REFRESH_EXPIRATION_DAYS": "30"
}
```

### Production Environment
**Secret Name:** `production/suitemedia/app`

```json
{
  "DB_HOST": "production-aurora.cluster-xxx.rds.amazonaws.com",
  "DB_PORT": "5432",
  "DB_USER": "suitemedia_user",
  "DB_PASSWORD": "production-strong-password",
  "DB_NAME": "suitemedia_production",
  "DB_SSLMODE": "require",
  "REDIS_HOST": "production-redis.cache.amazonaws.com",
  "REDIS_PORT": "6379",
  "REDIS_PASSWORD": "production-redis-password",
  "JWT_SECRET": "production-jwt-secret-strong-key",
  "JWT_REFRESH_SECRET": "production-refresh-secret-strong",
  "JWT_EXPIRATION_HOURS": "24",
  "JWT_REFRESH_EXPIRATION_DAYS": "30",
  "CORS_ALLOWED_ORIGINS": "https://suitemedia.com,https://www.suitemedia.com"
}
```

## Setup Instructions

### 1. Create AWS Secrets

**Staging:**
```bash
aws secretsmanager create-secret \
  --name staging/suitemedia/app \
  --description "Staging environment for SuiteMedia API" \
  --secret-string file://staging-secrets.json \
  --region ap-southeast-3
```

**Production:**
```bash
aws secretsmanager create-secret \
  --name production/suitemedia/app \
  --description "Production environment for SuiteMedia API" \
  --secret-string file://production-secrets.json \
  --region ap-southeast-3
```

### 2. Configure GitHub Secrets

Go to **GitHub Repository → Settings → Secrets and variables → Actions**

Add these repository secrets:
- `AWS_ACCESS_KEY_ID` - IAM user access key
- `AWS_SECRET_ACCESS_KEY` - IAM user secret key
- `SLACK_WEBHOOK_URL` - (Optional) For deployment notifications

### 3. Deployment Flow

#### Staging Deployment (Develop Branch)
```bash
# 1. Create feature branch
git checkout -b feature/new-feature

# 2. Make changes and commit
git add .
git commit -m "Add new feature"

# 3. Push to develop branch
git checkout develop
git merge feature/new-feature
git push origin develop

# 🚀 GitHub Actions automatically:
# - Builds and tests code
# - Creates Docker image
# - Pushes to ECR
# - Deploys to staging EKS cluster
# - Runs smoke tests
```

#### Production Deployment (Main Branch)
```bash
# 1. Merge develop to main after testing
git checkout main
git merge develop
git push origin main

# 🚀 GitHub Actions automatically:
# - Builds and tests code  
# - Creates Docker image
# - Pushes to ECR
# - Deploys to production EKS cluster
# - Runs smoke tests
# - Monitors for 5 minutes
# - Sends Slack notification
```

## How Secrets Are Used

### 1. GitHub Actions Workflow
```yaml
# GitHub Actions retrieves secrets for validation/configuration
- name: Get Secrets from AWS Secrets Manager
  run: |
    SECRET_JSON=$(aws secretsmanager get-secret-value \
      --secret-id staging/suitemedia/app \
      --query SecretString --output text)
    
    # Masks secrets in logs
    echo "::add-mask::$(echo $SECRET_JSON | jq -r '.DB_PASSWORD')"
```

### 2. External Secrets Operator (In Kubernetes)
The **External Secrets Operator** automatically syncs secrets from AWS to Kubernetes:

```yaml
# values-staging.yaml
externalSecrets:
  enabled: true
  data:
    - secretKey: DB_PASSWORD
      remoteRef:
        key: staging/suitemedia/app    # AWS Secret name
        property: DB_PASSWORD            # JSON property
```

### 3. Application Pods
Secrets are mounted as environment variables:

```yaml
# deployment.yaml
env:
  - name: DB_PASSWORD
    valueFrom:
      secretKeyRef:
        name: suitemedia-secrets      # Created by External Secrets
        key: DB_PASSWORD
```

## Verification Commands

### Check Staging Deployment
```bash
# Update kubeconfig
aws eks update-kubeconfig --name production-eks-cluster-staging --region ap-southeast-3

# Check pods
kubectl get pods -n staging -l app=suitemedia-app

# Check External Secrets
kubectl get externalsecrets -n staging
kubectl get secrets -n staging | grep suitemedia

# Check logs
kubectl logs -n staging deployment/suitemedia-app --tail=50

# Test health endpoint
curl https://staging.suitemedia.com/health
```

### Check Production Deployment
```bash
# Update kubeconfig
aws eks update-kubeconfig --name production-eks-cluster --region ap-southeast-3

# Check pods
kubectl get pods -n production -l app=suitemedia-app

# Check External Secrets
kubectl get externalsecrets -n production
kubectl describe externalsecret suitemedia-app -n production

# Check synced secrets
kubectl get secrets suitemedia-secrets -n production -o yaml

# Test health endpoint
curl https://suitemedia.com/health
curl https://suitemedia.com/ready
```

## Environment Differences

| Feature | Staging | Production |
|---------|---------|------------|
| **Replicas** | 2-5 pods | 6-20 pods |
| **Node Type** | Spot instances | On-demand |
| **Monitoring** | Basic | Full (Prometheus + Grafana) |
| **Backup** | Daily | Hourly + Daily |
| **Resources** | 256Mi-512Mi RAM | 512Mi-1Gi RAM |
| **Autoscaling** | CPU 80% trigger | CPU 70% trigger |
| **Health Check** | Every 10s | Every 5s |
| **Secret Refresh** | Every 30 min | Every 1 hour |

## Troubleshooting

### Secrets Not Syncing
```bash
# Check External Secrets Operator
kubectl get pods -n external-secrets-system

# Check SecretStore
kubectl get secretstore -n staging
kubectl describe secretstore aws-secrets-manager -n staging

# Check ExternalSecret status
kubectl describe externalsecret suitemedia-app -n staging
```

### Deployment Failed
```bash
# Check rollout status
kubectl rollout status deployment/suitemedia-app -n staging

# Check pod events
kubectl describe pod <pod-name> -n staging

# Check logs
kubectl logs -n staging deployment/suitemedia-app --previous
```

### Update Secrets
```bash
# Update secret in AWS
aws secretsmanager update-secret \
  --secret-id staging/suitemedia/app \
  --secret-string '{"DB_PASSWORD":"new-password"}'

# External Secrets Operator will auto-sync (wait 30min for staging, 1h for prod)
# Or force immediate sync by deleting the ExternalSecret
kubectl delete externalsecret suitemedia-app -n staging
kubectl apply -f helm/suitemedia/templates/external-secret.yaml
```

## Security Best Practices

✅ Different secrets for staging and production  
✅ Secrets never committed to Git  
✅ Secrets automatically masked in GitHub Actions logs  
✅ Secrets encrypted at rest in AWS Secrets Manager  
✅ IAM roles with least privilege  
✅ Automatic secret rotation (configure in AWS)  
✅ Audit logging enabled for secret access  
✅ Secrets synced automatically via External Secrets Operator

## Quick Reference

### Deploy to Staging
```bash
git push origin develop
```

### Deploy to Production
```bash
git push origin main
```

### Rollback
```bash
helm rollback suitemedia-app -n production
```

### View Deployment History
```bash
helm history suitemedia-app -n production
```
