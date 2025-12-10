# Helm Chart for SuiteMedia Application

Production-ready Helm chart for deploying SuiteMedia application on Kubernetes with full observability, auto-scaling, and security features.

## Features

✅ **High Availability**
- Multi-replica deployment with pod anti-affinity
- Pod Disruption Budget (PDB)
- Rolling updates with zero downtime
- Health checks (liveness, readiness, startup probes)

✅ **Auto-Scaling**
- Horizontal Pod Autoscaler (HPA) - CPU & Memory based
- Vertical Pod Autoscaler (VPA) - Optional
- Custom scaling behaviors (scale-up/scale-down policies)

✅ **Security**
- Non-root containers
- Read-only root filesystem
- Security contexts & Pod Security Standards
- Network policies (ingress/egress)
- External Secrets Operator integration with AWS Secrets Manager
- IAM Roles for Service Accounts (IRSA)

✅ **Observability**
- Prometheus metrics (ServiceMonitor)
- Grafana dashboards
- Structured logging (JSON)
- Filebeat log shipping
- Request tracing

✅ **Networking**
- Kong API Gateway integration
- Rate limiting & CORS
- TLS/SSL termination
- Multi-domain support

✅ **Disaster Recovery**
- Velero backup integration
- Automated backup schedules
- Point-in-time recovery

## Prerequisites

- Kubernetes 1.28+
- Helm 3.13+
- kubectl configured
- AWS EKS cluster (or any Kubernetes cluster)
- External Secrets Operator installed
- Prometheus Operator (optional, for ServiceMonitor)
- Kong Ingress Controller

## Installation

### 1. Add Helm Repository (if published)

```bash
helm repo add suitemedia https://charts.suitemedia.com
helm repo update
```

### 2. Create Namespace

```bash
kubectl create namespace production
```

### 3. Install External Secrets Operator

```bash
helm repo add external-secrets https://charts.external-secrets.io
helm install external-secrets \
  external-secrets/external-secrets \
  -n external-secrets-system \
  --create-namespace
```

### 4. Create SecretStore for AWS Secrets Manager

```yaml
# secret-store.yaml
apiVersion: external-secrets.io/v1beta1
kind: SecretStore
metadata:
  name: aws-secrets-manager
  namespace: production
spec:
  provider:
    aws:
      service: SecretsManager
      region: us-east-1
      auth:
        jwt:
          serviceAccountRef:
            name: suitemedia-production
```

```bash
kubectl apply -f secret-store.yaml
```

### 5. Install Chart - Staging

```bash
helm install suitemedia-staging ./helm/suitemedia \
  --namespace staging \
  --create-namespace \
  --values ./helm/suitemedia/values-staging.yaml \
  --set image.repository=123456789.dkr.ecr.us-east-1.amazonaws.com/suitemedia-app \
  --set image.tag=develop-abc1234 \
  --wait \
  --timeout 10m
```

### 6. Install Chart - Production

```bash
helm install suitemedia ./helm/suitemedia \
  --namespace production \
  --create-namespace \
  --values ./helm/suitemedia/values-production.yaml \
  --set image.repository=123456789.dkr.ecr.us-east-1.amazonaws.com/suitemedia-app \
  --set image.tag=v1.0.0 \
  --wait \
  --timeout 10m
```

## Upgrading

### Staging

```bash
helm upgrade suitemedia-staging ./helm/suitemedia \
  --namespace staging \
  --values ./helm/suitemedia/values-staging.yaml \
  --set image.tag=develop-xyz5678 \
  --wait \
  --timeout 10m
```

### Production

```bash
helm upgrade suitemedia ./helm/suitemedia \
  --namespace production \
  --values ./helm/suitemedia/values-production.yaml \
  --set image.tag=v1.1.0 \
  --wait \
  --timeout 10m
```

## Rollback

```bash
# List releases
helm history suitemedia -n production

# Rollback to previous version
helm rollback suitemedia -n production

# Rollback to specific revision
helm rollback suitemedia 3 -n production
```

## Configuration

### Key Configuration Values

| Parameter | Description | Default | Production |
|-----------|-------------|---------|------------|
| `replicaCount` | Number of pod replicas | `3` | `6` |
| `image.repository` | Container image repository | `""` | Set by CI/CD |
| `image.tag` | Container image tag | `""` | Set by CI/CD |
| `autoscaling.minReplicas` | Minimum HPA replicas | `3` | `6` |
| `autoscaling.maxReplicas` | Maximum HPA replicas | `10` | `20` |
| `resources.requests.cpu` | CPU request | `250m` | `500m` |
| `resources.requests.memory` | Memory request | `256Mi` | `512Mi` |
| `resources.limits.cpu` | CPU limit | `1000m` | `1000m` |
| `resources.limits.memory` | Memory limit | `1Gi` | `1Gi` |
| `ingress.enabled` | Enable ingress | `true` | `true` |
| `ingress.hosts` | Ingress hostnames | `suitemedia.com` | `suitemedia.com` |

### Override Values

Create custom values file:

```yaml
# custom-values.yaml
replicaCount: 4

resources:
  requests:
    cpu: 300m
    memory: 384Mi
  limits:
    cpu: 1500m
    memory: 1.5Gi

autoscaling:
  minReplicas: 4
  maxReplicas: 15
  targetCPUUtilizationPercentage: 75

ingress:
  hosts:
    - host: app.custom-domain.com
      paths:
        - path: /
          pathType: Prefix
```

Install with custom values:

```bash
helm install suitemedia ./helm/suitemedia \
  -f custom-values.yaml \
  --namespace production
```

## Verification

### 1. Check Deployment Status

```bash
# Get all resources
kubectl get all -n production -l app.kubernetes.io/name=suitemedia

# Check deployment
kubectl get deployment -n production

# Check pods
kubectl get pods -n production

# Check HPA status
kubectl get hpa -n production

# Check ingress
kubectl get ingress -n production
```

### 2. View Logs

```bash
# All pods
kubectl logs -n production -l app.kubernetes.io/name=suitemedia --tail=100 -f

# Specific pod
kubectl logs -n production <pod-name> -f
```

### 3. Run Helm Tests

```bash
helm test suitemedia -n production
```

### 4. Check Metrics

```bash
# Pod metrics
kubectl top pods -n production -l app.kubernetes.io/name=suitemedia

# Node metrics
kubectl top nodes
```

### 5. Access Application

```bash
# Port forward (for testing)
kubectl port-forward -n production svc/suitemedia 8080:80

# Access via ingress
curl https://suitemedia.com/health
```

## Monitoring

### Prometheus Metrics

The application exposes metrics at `/metrics` endpoint. ServiceMonitor is automatically created if Prometheus Operator is installed.

```bash
# Check ServiceMonitor
kubectl get servicemonitor -n production
```

### Grafana Dashboards

Import Grafana dashboard from `monitoring/grafana-dashboard.json`

### Alerts

Prometheus alerts are configured in `values-production.yaml`:
- High error rate (> 5%)
- High latency (P95 > 1s)
- Pod restart count
- HPA maxed out

## Troubleshooting

### Pods Not Starting

```bash
# Describe pod
kubectl describe pod <pod-name> -n production

# Check events
kubectl get events -n production --sort-by='.lastTimestamp'

# Check logs
kubectl logs <pod-name> -n production --previous
```

### Image Pull Errors

```bash
# Verify image exists in ECR
aws ecr describe-images \
  --repository-name suitemedia-app \
  --image-ids imageTag=v1.0.0 \
  --region us-east-1

# Check ImagePullSecrets
kubectl get serviceaccount suitemedia-production -n production -o yaml
```

### External Secrets Not Working

```bash
# Check ExternalSecret status
kubectl get externalsecret -n production
kubectl describe externalsecret suitemedia -n production

# Check SecretStore
kubectl get secretstore -n production
kubectl describe secretstore aws-secrets-manager -n production

# Verify IAM permissions
aws iam get-role --role-name suitemedia-production-role
```

### HPA Not Scaling

```bash
# Check HPA status
kubectl describe hpa suitemedia -n production

# Check metrics server
kubectl get apiservice v1beta1.metrics.k8s.io -o yaml

# Check metrics availability
kubectl top pods -n production
```

### Ingress Issues

```bash
# Check ingress
kubectl describe ingress suitemedia -n production

# Check Kong service
kubectl get svc -n kong

# Test from within cluster
kubectl run -it --rm debug \
  --image=curlimages/curl \
  --restart=Never \
  -- curl http://suitemedia.production.svc.cluster.local/health
```

## Backup & Restore

### Velero Backup

```bash
# Create on-demand backup
velero backup create suitemedia-backup-$(date +%Y%m%d) \
  --include-namespaces production \
  --wait

# List backups
velero backup get

# Restore from backup
velero restore create --from-backup suitemedia-backup-20251210
```

### Database Backup

Database backups are handled separately via Aurora automated backups (configured in infrastructure).

## Uninstall

```bash
# Uninstall Helm release
helm uninstall suitemedia -n production

# Delete namespace (optional)
kubectl delete namespace production
```

## CI/CD Integration

This Helm chart is designed to work with the GitHub Actions CI/CD pipeline in `.github/workflows/ci-cd-pipeline.yml`.

### Deployment Flow

```
1. Code push to main/develop
2. Run tests (unit, integration, e2e)
3. Build Docker image
4. Push to ECR
5. Deploy with Helm:
   - develop → staging
   - main → production
```

### GitHub Actions Usage

```yaml
- name: Deploy to Production
  run: |
    helm upgrade --install suitemedia ./helm/suitemedia \
      --namespace production \
      --create-namespace \
      --set image.repository=$ECR_REGISTRY/$ECR_REPOSITORY \
      --set image.tag=$IMAGE_TAG \
      --values ./helm/suitemedia/values-production.yaml \
      --wait \
      --timeout 10m
```

## Best Practices

1. **Always use specific image tags** - Avoid `latest` tag in production
2. **Test in staging first** - Deploy to staging before production
3. **Monitor deployments** - Watch metrics during and after deployment
4. **Use Helm diffs** - Review changes before applying
5. **Keep values files in version control**
6. **Use External Secrets** - Never commit secrets to Git
7. **Enable PodDisruptionBudget** - Ensure availability during maintenance
8. **Set resource limits** - Prevent resource exhaustion
9. **Use health checks** - Configure proper liveness/readiness probes
10. **Regular backups** - Schedule automated Velero backups

## Support

For issues or questions:
- Check documentation: `docs/`
- View logs: `kubectl logs -n production`
- Contact: devops@suitemedia.com

## License

MIT License

---

**Chart Version:** 1.0.0  
**App Version:** 1.0.0  
**Maintained by:** DevOps Team  
**Last Updated:** December 10, 2025
