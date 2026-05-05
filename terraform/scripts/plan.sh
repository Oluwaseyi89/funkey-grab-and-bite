#!/usr/bin/env bash
# scripts/plan.sh — Preview changes for a specific component or full environment
# Usage:
#   ./scripts/plan.sh [production|staging] [component]
#
# Components match those in deploy.sh

set -euo pipefail

ENV="${1:-production}"
COMPONENT="${2:-all}"
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ENV_DIR="${DIR}/environments/${ENV}"

if [[ ! -d "$ENV_DIR" ]]; then
  echo "❌  Unknown environment: ${ENV}"
  exit 1
fi

cd "$ENV_DIR"

echo "🔄  terraform init..."
terraform init -upgrade -reconfigure

case "$COMPONENT" in
  infrastructure)
    terraform plan \
      -target=module.networking \
      -target=module.waf \
      -target=module.dns \
      -target=module.ecr \
      -target=module.secrets
    ;;
  database)   terraform plan -target=module.aurora ;;
  web)
    terraform plan \
      -target=module.static_hosting_web \
      -target=module.cloudfront_web
    ;;
  admin)
    terraform plan \
      -target=module.static_hosting_admin \
      -target=module.cloudfront_admin
    ;;
  api)
    terraform plan \
      -target=module.alb \
      -target=module.cloudfront_api \
      -target=module.ecs_api
    ;;
  workers)   terraform plan -target=module.ecs_workers ;;
  lambda)    terraform plan -target=module.lambda_catering ;;
  cache)     terraform plan -target=module.elasticache ;;
  queue)
    terraform plan \
      -target=module.sqs \
      -target=module.eventbridge
    ;;
  monitoring) terraform plan -target=module.monitoring ;;
  all)        terraform plan ;;
  *)
    echo "❌  Unknown component: ${COMPONENT}"
    exit 1
    ;;
esac
