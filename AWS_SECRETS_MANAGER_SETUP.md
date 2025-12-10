# AWS Secrets Manager Setup Guide

## Create Secrets in AWS Secrets Manager

### For Staging Environment

```bash
aws secretsmanager create-secret \
  --name staging/suitemedia/app \
  --description "Staging environment secrets for SuiteMedia API" \
  --secret-string '{
    "DB_HOST": "staging-aurora-cluster.cluster-xxx.ap-southeast-3.rds.amazonaws.com",
    "DB_PORT": "5432",
    "DB_USER": "suitemedia_user",
    "DB_PASSWORD": "your-staging-db-password",
    "DB_NAME": "suitemedia_staging",
    "DB_SSLMODE": "require",
    "REDIS_HOST": "staging-redis.xxx.cache.amazonaws.com",
    "REDIS_PORT": "6379",
    "REDIS_PASSWORD": "your-redis-password",
    "JWT_SECRET": "your-staging-jwt-secret-key",
    "JWT_REFRESH_SECRET": "your-staging-refresh-secret-key",
    "JWT_EXPIRATION_HOURS": "24",
    "JWT_REFRESH_EXPIRATION_DAYS": "30"
  }' \
  --region ap-southeast-3
```

### For Production Environment

```bash
aws secretsmanager create-secret \
  --name production/suitemedia/app \
  --description "Production environment secrets for SuiteMedia API" \
  --secret-string '{
    "DB_HOST": "production-aurora-cluster.cluster-xxx.ap-southeast-3.rds.amazonaws.com",
    "DB_PORT": "5432",
    "DB_USER": "suitemedia_user",
    "DB_PASSWORD": "your-production-db-password",
    "DB_NAME": "suitemedia_production",
    "DB_SSLMODE": "require",
    "REDIS_HOST": "production-redis.xxx.cache.amazonaws.com",
    "REDIS_PORT": "6379",
    "REDIS_PASSWORD": "your-redis-password",
    "JWT_SECRET": "your-production-jwt-secret-key",
    "JWT_REFRESH_SECRET": "your-production-refresh-secret-key",
    "JWT_EXPIRATION_HOURS": "24",
    "JWT_REFRESH_EXPIRATION_DAYS": "30",
    "CORS_ALLOWED_ORIGINS": "https://suitemedia.com,https://www.suitemedia.com"
  }' \
  --region ap-southeast-3
```

## Update Secrets

```bash
aws secretsmanager update-secret \
  --secret-id production/suitemedia/app \
  --secret-string '{
    "DB_PASSWORD": "new-password",
    "JWT_SECRET": "new-secret"
  }' \
  --region ap-southeast-3
```

## Retrieve Secrets (for testing)

```bash
aws secretsmanager get-secret-value \
  --secret-id production/suitemedia/app \
  --region ap-southeast-3 \
  --query SecretString \
  --output text | jq .
```

## GitHub Secrets Configuration

Add these secrets to your GitHub repository:

```
AWS_ACCESS_KEY_ID
AWS_SECRET_ACCESS_KEY
SLACK_WEBHOOK_URL (optional)
```

Go to: **Settings → Secrets and variables → Actions → New repository secret**

## IAM Policy for GitHub Actions

Create IAM user with this policy:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ecr:GetAuthorizationToken",
        "ecr:BatchCheckLayerAvailability",
        "ecr:GetDownloadUrlForLayer",
        "ecr:BatchGetImage",
        "ecr:PutImage",
        "ecr:InitiateLayerUpload",
        "ecr:UploadLayerPart",
        "ecr:CompleteLayerUpload"
      ],
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "eks:DescribeCluster",
        "eks:ListClusters"
      ],
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "secretsmanager:GetSecretValue",
        "secretsmanager:DescribeSecret"
      ],
      "Resource": [
        "arn:aws:secretsmanager:ap-southeast-3:*:secret:staging/suitemedia/*",
        "arn:aws:secretsmanager:ap-southeast-3:*:secret:production/suitemedia/*"
      ]
    }
  ]
}
```

## Kubernetes External Secrets Operator Setup

The External Secrets Operator is already configured in the Helm chart. It will automatically sync secrets from AWS Secrets Manager to Kubernetes secrets.

### Verify External Secrets

```bash
# Check if External Secrets Operator is installed
kubectl get pods -n external-secrets-system

# Check ExternalSecret resources
kubectl get externalsecrets -n production

# Check synced secrets
kubectl get secrets -n production | grep suitemedia
```

## Workflow Summary

1. **GitHub Actions** retrieves secrets from AWS Secrets Manager for deployment configuration
2. **External Secrets Operator** in Kubernetes syncs secrets automatically
3. **Application pods** mount secrets as environment variables
4. **Secrets rotation** is handled by AWS Secrets Manager with automatic updates

## Security Best Practices

✅ Never commit secrets to Git  
✅ Use IAM roles with least privilege  
✅ Enable secret rotation in AWS Secrets Manager  
✅ Use different secrets for staging and production  
✅ Mask secrets in CI/CD logs  
✅ Regularly audit secret access  
✅ Use KMS encryption for secrets at rest
