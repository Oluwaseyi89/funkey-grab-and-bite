#!/usr/bin/env bash
# scripts/deploy.sh — Deploy a specific component or full environment
# Usage:
#   ./scripts/deploy.sh [production|staging] [component]
#
# Available components:
#   infrastructure  networking + WAF + DNS + ECR + secrets
#   database        Aurora PostgreSQL
#   web             web app (S3 + CloudFront)
#   admin           admin app (S3 + CloudFront)
#   api             ALB + ECS API + CloudFront API
#   workers         ECS worker cluster
#   lambda          Lambda catering-notify
#   cache           ElastiCache Redis
#   queue           SQS + EventBridge
#   monitoring      CloudWatch + SNS
#   all             full environment (no targets)

set -euo pipefail

ENV="${1:-production}"
COMPONENT="${2:-all}"
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ENV_DIR="${DIR}/environments/${ENV}"

if [[ ! -d "$ENV_DIR" ]]; then
  echo "❌  Unknown environment: ${ENV}  (use 'production' or 'staging')"
  exit 1
fi

cd "$ENV_DIR"

if [[ ! -f "terraform.tfvars" ]]; then
  echo "⚠️   terraform.tfvars not found — copying from example"
  cp terraform.tfvars.example terraform.tfvars
  echo "✏️   Edit ${ENV_DIR}/terraform.tfvars with real values then re-run"
  exit 1
fi

echo "🔄  terraform init..."
terraform init -upgrade -reconfigure

case "$COMPONENT" in
  infrastructure)
    echo "🏗️   Deploying core infrastructure (networking + WAF + DNS + ECR + secrets)..."
    terraform apply \
      -target=module.networking \
      -target=module.waf \
      -target=module.dns \
      -target=module.ecr \
      -target=module.secrets \
      -auto-approve
    ;;
  database)
    echo "🗃️   Deploying Aurora PostgreSQL..."
    terraform apply -target=module.aurora -auto-approve
    ;;
  web)
    echo "🌐  Deploying web app (S3 + CloudFront)..."
    terraform apply \
      -target=module.static_hosting_web \
      -target=module.cloudfront_web \
      -auto-approve
    ;;
  admin)
    echo "🔧  Deploying admin app (S3 + CloudFront)..."
    terraform apply \
      -target=module.static_hosting_admin \
      -target=module.cloudfront_admin \
      -auto-approve
    ;;
  api)
    echo "⚙️   Deploying Go API (ALB + ECS + CloudFront edge)..."
    terraform apply \
      -target=module.alb \
      -target=module.cloudfront_api \
      -target=module.ecs_api \
      -auto-approve
    ;;
  workers)
    echo "👷  Deploying ECS worker cluster..."
    terraform apply -target=module.ecs_workers -auto-approve
    ;;
  lambda)
    echo "λ   Deploying Lambda catering-notify..."
    terraform apply -target=module.lambda_catering -auto-approve
    ;;
  cache)
    echo "⚡  Deploying ElastiCache Redis..."
    terraform apply -target=module.elasticache -auto-approve
    ;;
  queue)
    echo "📬  Deploying SQS queues + EventBridge scheduler..."
    terraform apply \
      -target=module.sqs \
      -target=module.eventbridge \
      -auto-approve
    ;;
  monitoring)
    echo "📊  Deploying monitoring (CloudWatch + SNS)..."
    terraform apply -target=module.monitoring -auto-approve
    ;;
  all)
    echo "🚀  Deploying full ${ENV} environment..."
    terraform apply -auto-approve
    ;;
  *)
    echo "❌  Unknown component: ${COMPONENT}"
    echo "    Valid components: infrastructure, database, web, admin, api, workers, lambda, cache, queue, monitoring, all"
    exit 1
    ;;
esac

echo "✅  Done! Run 'terraform -chdir=${ENV_DIR} output' to see outputs."
