# Funkey Grab & Bite — **AWS Infrastructure (Terraform)**

[![Terraform](https://img.shields.io/badge/Terraform-%3E%3D1.6-844FBA?logo=terraform)](https://www.terraform.io/)
[![AWS Provider](https://img.shields.io/badge/AWS%20Provider-~%3E5.50-FF9900?logo=amazon-aws)](https://registry.terraform.io/providers/hashicorp/aws/latest)
[![Region](https://img.shields.io/badge/Region-us--east--1-232F3E?logo=amazon-aws)](https://aws.amazon.com/about-aws/global-infrastructure/regions_az/)
[![License](https://img.shields.io/badge/License-MIT-green)](../LICENSE)

Production-grade, modular AWS infrastructure for the **Funkey Grab & Bite** food-ordering platform. Every component is an independently deployable Terraform module controlled by a feature flag — perfect for portfolio showcasing a single service or spinning up the full stack in one command.

---

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Directory Structure](#directory-structure)
3. [Module Reference](#module-reference)
4. [Prerequisites](#prerequisites)
5. [First-Time Setup](#first-time-setup)
6. [Preflight Checks (No Apply)](#preflight-checks-no-apply)
7. [Feature Flags](#feature-flags)
8. [Deploy Order (Full Stack)](#deploy-order-full-stack)
9. [Targeted Deploys](#targeted-deploys)
10. [Environment Differences](#environment-differences)
11. [Post-Deploy Steps](#post-deploy-steps)
12. [CI/CD Integration](#cicd-integration)
13. [Makefile Reference](#makefile-reference)
14. [Security Controls](#security-controls)
15. [Troubleshooting](#troubleshooting)

---

## Architecture Overview

```
                        ┌──────────────────────────────────────────────────────┐
                        │                   Route 53 (DNS)                      │
                        │  funkeygrabandbite.com  /  api.  /  admin.  /  www.  │
                        └──────────────┬───────────────────────────────────────┘
                                       │
                        ┌──────────────▼───────────────────────────────────────┐
                        │           AWS WAF v2 + Shield Standard               │
                        │   CRS · KnownBadInputs · SQLi · Rate-limit (2k/5m)  │
                        └──────┬───────────────────────────┬────────────────────┘
                               │                           │
              ┌────────────────▼──────┐       ┌───────────▼──────────────────┐
              │  CloudFront (Web SPA) │       │  CloudFront (Admin SPA)      │
              │  Nuxt 4 Storefront    │       │  React 19 Dashboard          │
              └────────┬─────────────┘       └───────────┬──────────────────┘
                       │ OAC                              │ OAC
              ┌────────▼─────────────┐       ┌───────────▼──────────────────┐
              │  S3 Bucket (Web)     │       │  S3 Bucket (Admin)           │
              │  Private + Versioned │       │  Private + Versioned         │
              └──────────────────────┘       └──────────────────────────────┘

                        ┌──────────────────────────────────────┐
                        │       CloudFront (API Edge)          │
                        │  CachingDisabled · AllViewer RP      │
                        │  Origin Shield · WebSocket support   │
                        └────────────────┬─────────────────────┘
                                         │ secret header
                        ┌────────────────▼─────────────────────┐
                        │     Application Load Balancer         │
                        │  HTTP→HTTPS redirect · Header check   │
                        └────────────────┬─────────────────────┘
                                         │
          ┌──────────────────────────────┴────────────────────────────────┐
          │                                                               │
┌─────────▼──────────────┐                               ┌───────────────▼───────────┐
│   ECS Fargate (API)    │                               │  ECS Fargate (Workers)    │
│   Go/Gin REST API      │                               │  Background task runners  │
│   port 8080            │                               │  FARGATE_SPOT preferred   │
│   2 tasks (prod)       │                               │  SQS consumer             │
└─────────┬──────────────┘                               └───────────────────────────┘
          │
┌─────────▼────────────────────────────────────────────────────────────────────────┐
│                               Data / Messaging Layer                              │
│                                                                                   │
│  ┌──────────────────────────┐  ┌────────────────────┐  ┌──────────────────────┐  │
│  │  Aurora PostgreSQL Sv2   │  │  ElastiCache Redis  │  │  SQS (Order + Catr.) │  │
│  │  Multi-AZ · Auto-scale   │  │  7.1 · allkeys-lru  │  │  + DLQs · HTTPS-only│  │
│  └──────────────────────────┘  └────────────────────┘  └──────────────────────┘  │
└──────────────────────────────────────────────────────────────────────────────────┘

┌───────────────────────────────────────────────────────────────────────────────────┐
│                            Async / Scheduled Work                                 │
│                                                                                   │
│  ┌──────────────────────────────────┐        ┌──────────────────────────────────┐ │
│  │   Lambda (Catering Notifier)     │        │   EventBridge Scheduler          │ │
│  │   Container image · VPC          │        │   Promo expiry  00:05 UTC daily  │ │
│  │   SQS trigger · batch 5          │        │   Reports gen   01:00 UTC daily  │ │
│  └──────────────────────────────────┘        └──────────────────────────────────┘ │
└───────────────────────────────────────────────────────────────────────────────────┘

┌───────────────────────────────────────────────────────────────────────────────────┐
│                              Security & Secrets                                   │
│   Secrets Manager (JWT · SES · Twilio · Paystack)  ·  IAM Roles  ·  VPC SGs      │
└───────────────────────────────────────────────────────────────────────────────────┘

┌───────────────────────────────────────────────────────────────────────────────────┐
│                                Observability                                      │
│   CloudWatch Logs · CloudWatch Alarms (7) · Dashboard · SNS Alerts · X-Ray       │
└───────────────────────────────────────────────────────────────────────────────────┘
```

---

## Directory Structure

```
terraform/
├── Makefile                        # Convenience make targets (see §Makefile Reference)
├── .gitignore                      # Excludes state, secrets, plan files
├── scripts/
│   ├── bootstrap-state.sh          # One-time: create S3 bucket + DynamoDB table
│   ├── deploy.sh                   # Targeted apply (ENV + component args)
│   ├── plan.sh                     # Targeted plan  (ENV + component args)
│   └── tf-all.sh                   # Safe fmt/validate/plan runner for both envs
│
├── modules/                        # Reusable, independently callable modules
│   ├── networking/                 # VPC, subnets, NAT, SGs, VPC endpoints
│   ├── waf/                        # WAF v2 WebACL (CLOUDFRONT scope)
│   ├── dns/                        # Route 53 zone, ACM wildcard cert, health check
│   ├── static-hosting/             # Private S3 bucket + OAC (reusable for web/admin)
│   ├── cloudfront-spa/             # CloudFront for SPA apps (reusable for web/admin)
│   ├── cloudfront-api/             # CloudFront for ALB origin (API edge)
│   ├── alb/                        # Application Load Balancer + secret header guard
│   ├── ecr/                        # Two ECR repositories (api + lambda-catering)
│   ├── ecs-api/                    # ECS Fargate cluster/service for the Go API
│   ├── ecs-workers/                # ECS Fargate cluster/service for background workers
│   ├── lambda-catering/            # Lambda container (catering notifications + SQS)
│   ├── aurora/                     # Aurora PostgreSQL Serverless v2 (Multi-AZ)
│   ├── elasticache/                # ElastiCache Redis 7.1 replication group
│   ├── sqs/                        # Order + catering queues with DLQs
│   ├── eventbridge/                # Custom event bus + Scheduler rules
│   ├── secrets/                    # Secrets Manager secrets (app, paystack, admin)
│   └── monitoring/                 # CloudWatch alarms, dashboard, SNS
│
└── environments/
    ├── production/                 # Root config — 3 AZs, full HA, protect on
    │   ├── backend.tf              # Remote state: S3 + DynamoDB lock
    │   ├── versions.tf             # Provider versions + default_tags
    │   ├── variables.tf            # Feature flags + tunable params
    │   ├── locals.tf               # try() cross-module refs + common_tags
    │   ├── main.tf                 # All 19 module calls (count = feature flag)
    │   ├── outputs.tf              # URLs, ARNs, endpoints
    │   └── terraform.tfvars.example
    └── staging/                    # Root config — 2 AZs, single NAT, cost-optimised
        └── (same structure)
```

---

## Module Reference

| Module | Purpose | Key Outputs |
|---|---|---|
| `networking` | VPC (`10.0.0.0/16`), 3-tier subnets, NAT GW, 6 SGs, VPC endpoints | `vpc_id`, `*_subnet_ids`, `*_sg_id` |
| `waf` | WAF v2 WebACL (CRS + SQLi + rate-limit) | `web_acl_arn` |
| `dns` | Route 53 zone, ACM wildcard cert, API health check | `zone_id`, `certificate_arn` |
| `static-hosting` | Private versioned S3 + OAC (reused for web & admin) | `bucket_id`, `oac_id` |
| `cloudfront-spa` | CloudFront SPA distribution with 404→`index.html` | `distribution_id`, `cloudfront_domain_name` |
| `cloudfront-api` | CloudFront → ALB with secret header, Origin Shield | `cloudfront_domain_name`, `origin_secret_*` |
| `alb` | Internet-facing ALB, HTTP→HTTPS, secret header guard | `alb_dns_name`, `target_group_arn` |
| `ecr` | Two private repositories with lifecycle policies | `api_repository_url`, `lambda_catering_repository_url` |
| `ecs-api` | Fargate cluster/service, Secrets Manager env injection, auto-scaling | `cluster_arn`, `service_name` |
| `ecs-workers` | Fargate (SPOT preferred) workers cluster, SQS consumer | `cluster_arn`, `task_role_arn` |
| `lambda-catering` | Container Lambda, SQS trigger (batch 5), VPC placement | `function_arn`, `function_name` |
| `aurora` | Aurora PostgreSQL Serverless v2, writer + optional reader | `cluster_endpoint`, `db_secret_arn` |
| `elasticache` | Redis 7.1 replication group, Multi-AZ, allkeys-lru | `primary_endpoint_address` |
| `sqs` | Order queue + catering queue, each with DLQ, HTTPS-only policy | `order_queue_url`, `catering_queue_url` |
| `eventbridge` | Custom event bus + Scheduler (promo expiry + daily report) | `event_bus_arn` |
| `secrets` | Secrets Manager — app bundle, Paystack, admin bootstrap | `app_secrets_arn`, `paystack_secret_arn` |
| `monitoring` | SNS email alerts, 7 CloudWatch alarms, dashboard | `sns_topic_arn`, `dashboard_name` |

---

## Prerequisites

| Tool | Version | Purpose |
|---|---|---|
| [Terraform](https://developer.hashicorp.com/terraform/install) | `>= 1.6.0` | Infrastructure as code runner |
| [AWS CLI](https://aws.amazon.com/cli/) | `>= 2.x` | Bootstrap script + ECR login |
| [Docker](https://docs.docker.com/engine/install/) | any | Build + push container images |
| AWS credentials | — | Configured via `~/.aws/credentials` or env vars |

**Minimum IAM permissions** for the deploying principal:

- `AdministratorAccess` is simplest for initial setup.
- For least-privilege, the principal needs: `ec2:*`, `ecs:*`, `ecr:*`, `rds:*`, `elasticache:*`, `lambda:*`, `s3:*`, `cloudfront:*`, `route53:*`, `acm:*`, `wafv2:*`, `sqs:*`, `events:*`, `scheduler:*`, `secretsmanager:*`, `cloudwatch:*`, `sns:*`, `iam:*` (for role creation), `logs:*`.

---

## First-Time Setup

### 1. Bootstrap remote state (run once per AWS account)

```bash
cd terraform/
./scripts/bootstrap-state.sh
```

This creates:
- S3 bucket `funkey-terraform-state` (versioned, AES-256, private)
- DynamoDB table `funkey-terraform-locks` (PAY_PER_REQUEST)

If your AWS account already has a different state bucket/region, update both backend files before init:

- `environments/staging/backend.tf`
- `environments/production/backend.tf`

### 2. Configure your variables

```bash
# Production
cp environments/production/terraform.tfvars.example environments/production/terraform.tfvars
nano environments/production/terraform.tfvars

# Staging
cp environments/staging/terraform.tfvars.example environments/staging/terraform.tfvars
nano environments/staging/terraform.tfvars
```

**Required values to fill in** (no defaults):

| Variable | Description |
|---|---|
| `alert_email` | Email address for CloudWatch alarm notifications |
| `ses_sender_email` | Verified SES email address for the API |
| `api_image_tag` | ECR image tag for the Go API (e.g. `latest`) |

Note: Lambda currently reuses `api_image_tag`.

### 3. Initialise Terraform

```bash
make init ENV=production
make init ENV=staging
```

### 4. Review the plan

```bash
make plan ENV=staging
```

### 5. Deploy (see §Deploy Order below)

---

## Preflight Checks (No Apply)

Run these checks before any apply:

```bash
cd terraform

# Formatting
./scripts/tf-all.sh fmt

# Validate both environments without backend state access
./scripts/tf-all.sh validate

# Plan both environments in a temporary backend-free workspace
./scripts/tf-all.sh plan
```

What this gives you:
- catches HCL/schema errors early
- avoids accidental applies
- avoids backend lock/state issues during quick verification

---

## Feature Flags

Each module is guarded by a boolean flag inside the `features` variable. Set flags to `false` to skip components you don't need for a particular showcase.

```hcl
# environments/production/terraform.tfvars
features = {
  networking      = true    # VPC, subnets, SG, endpoints
  waf             = true    # WAF v2 WebACL (CloudFront)
  dns             = true    # Route53 zone + ACM + records
  web_app         = true    # CloudFront + S3 for Nuxt web app
  admin_app       = true    # CloudFront + S3 for React admin
  api             = true    # ALB + API edge + ECS API stack
  workers         = true    # Background worker tasks on ECS
  lambda_catering  = true    # Catering notifier Lambda
  database         = true    # Aurora PostgreSQL Serverless v2
  cache            = true    # ElastiCache Redis
  queue            = true    # SQS order + catering queues
  scheduler        = true    # EventBridge Scheduler rules
  secrets          = true    # Secrets Manager entries
  monitoring       = true    # CloudWatch alarms, dashboard, SNS
}
```

### Showcase Presets

| Goal | Flags to enable |
|---|---|
| **Minimal API** | `networking` + `dns` + `api` + `database` + `secrets` |
| **Frontend only** | `networking` + `dns` + `web_app` + `admin_app` (+ optional `waf`) |
| **Full backend** | `networking` + `dns` + `api` + `workers` + `database` + `cache` + `queue` + `scheduler` + `secrets` + `monitoring` |
| **Event-driven demo** | `networking` + `dns` + `api` + `workers` + `queue` + `lambda_catering` + `scheduler` + `secrets` |
| **Full production** | all `true` (default) |

---

## Deploy Order (Full Stack)

When deploying a brand-new environment, follow this order to satisfy cross-module dependencies:

```
Step 1 — Foundation
  make deploy-infrastructure ENV=production
  # Creates: VPC, subnets, SGs, WAF, Route 53 zone, ACM cert

Step 2 — Repositories (before images can be pushed)
  make deploy ENV=production   # or individual:
  terraform -chdir=environments/production apply -target=module.ecr

Step 3 — Data layer
  make deploy-database ENV=production
  # Creates: Aurora cluster + Secrets Manager secrets

Step 4 — Compute
  make deploy-api ENV=production
  # Creates: ALB → ECS Fargate (API) + CloudFront edge

Step 5 — Web & Admin frontends
  make deploy-web   ENV=production
  make deploy-admin ENV=production

Step 6 — Async / messaging
  make deploy-queue     ENV=production
  make deploy-workers   ENV=production
  make deploy-lambda    ENV=production
  make deploy-scheduler ENV=production

Step 7 — Caching
  make deploy-cache ENV=production

Step 8 — Observability
  make deploy-monitoring ENV=production
```

Or deploy everything at once (safe with feature flags):

```bash
make deploy ENV=production
```

---

## Targeted Deploys

### Using the Makefile

```bash
# Deploy only the web frontend (e.g. after a Nuxt build)
make deploy-web ENV=staging

# Plan only the API changes before applying
make plan-api ENV=production

# Deploy only the monitoring stack
make deploy-monitoring ENV=production
```

### Using Terraform -target directly

```bash
cd environments/production

# Update a single module
terraform apply -target=module.ecs_api

# Update two related modules together
terraform apply -target=module.aurora -target=module.secrets

# Destroy only the Lambda (for cost saving)
terraform destroy -target=module.lambda_catering
```

### Using the deploy script directly

```bash
# deploy.sh <env> <component>
./scripts/deploy.sh staging web
./scripts/deploy.sh production api
./scripts/plan.sh   staging   database
```

---

## Environment Differences

| Setting | Production | Staging |
|---|---|---|
| AZ count | 3 | 2 |
| NAT Gateways | 1 per AZ (3 total) | 1 shared |
| Aurora min capacity | 0.5 ACU | 0 ACU (scales to zero) |
| Aurora max capacity | 8 ACU | 2 ACU |
| Aurora reader | enabled | disabled |
| ECS desired count | 2 | 1 |
| ECS task CPU | 512 | 256 |
| ECS task memory | 1024 MB | 512 MB |
| ElastiCache nodes | 2 (Multi-AZ) | 1 |
| ElastiCache type | `cache.t4g.small` | `cache.t4g.micro` |
| WAF enabled | true | false (save cost) |
| Workers enabled | true | false |
| Lambda catering | true | false |
| Scheduler | true | false |
| S3 force_destroy | false | true |
| Secrets recovery window | 30 days | 0 (instant) |
| Deletion protection (Aurora) | true | false |

---

## Post-Deploy Steps

After a successful `make deploy ENV=production`, complete these manual steps:

### 1. Update Domain Nameservers

```bash
# Get the Route 53 nameservers
make output ENV=production | grep name_servers
```

Copy the 4 nameserver values into your domain registrar's NS records for `funkeygrabandbite.com`.

> ACM certificate validation is automated via DNS — it will complete once NS propagation finishes (up to 48 hours for new domains).

### 2. Push Container Images to ECR

```bash
# Get ECR URLs
make output-json ENV=production | jq '.api_ecr_url.value,.lambda_ecr_url.value'

# Authenticate Docker with ECR
aws ecr get-login-password --region us-east-1 \
  | docker login --username AWS --password-stdin <ACCOUNT_ID>.dkr.ecr.us-east-1.amazonaws.com

# Build and push the Go API
cd funkey-bite-api/
docker build -t funkey-bite-api:latest .
docker tag  funkey-bite-api:latest <API_ECR_URL>:latest
docker push <API_ECR_URL>:latest

# Build and push the Lambda image
docker build -t funkey-bite-lambda-catering:latest -f Dockerfile.lambda .
docker tag  funkey-bite-lambda-catering:latest <LAMBDA_ECR_URL>:latest
docker push <LAMBDA_ECR_URL>:latest
```

### 3. Populate Secrets Manager

Terraform creates the secret shells with placeholder values. Update them via the AWS Console or CLI:

```bash
# App secrets bundle (JWT, SES, Twilio)
aws secretsmanager put-secret-value \
  --secret-id funkey-production-app-secrets \
  --secret-string '{
    "JWT_SECRET":          "your-jwt-secret",
    "SES_SENDER_EMAIL":    "noreply@funkeygrabandbite.com",
    "TWILIO_ACCOUNT_SID":  "ACxxx",
    "TWILIO_AUTH_TOKEN":   "xxx",
    "TWILIO_FROM_NUMBER":  "+1xxx"
  }'

# Paystack credentials
aws secretsmanager put-secret-value \
  --secret-id funkey-production-paystack \
  --secret-string '{
    "PAYSTACK_SECRET_KEY": "sk_live_xxx",
    "PAYSTACK_PUBLIC_KEY": "pk_live_xxx"
  }'
```

> Terraform has `lifecycle { ignore_changes = [secret_string] }` on all secrets, so re-running `terraform apply` will **never** overwrite your live credentials.

### 4. Confirm SNS Email Subscription

Check your `alert_email` inbox and click the confirmation link from AWS SNS to activate CloudWatch alarm notifications.

### 5. Upload Static Assets to S3

```bash
# Web (Nuxt)
cd funkey-bite-web/
npm run generate
aws s3 sync .output/public/ s3://<WEB_BUCKET_ID>/ --delete

# Admin (React/Vite)
cd funkey-bite-admin/
npm run build
aws s3 sync dist/ s3://<ADMIN_BUCKET_ID>/ --delete

# Invalidate CloudFront caches
aws cloudfront create-invalidation --distribution-id <WEB_CF_ID>   --paths "/*"
aws cloudfront create-invalidation --distribution-id <ADMIN_CF_ID> --paths "/*"
```

---

## CI/CD Integration

### GitHub Actions — API Deploy

```yaml
# .github/workflows/deploy-api.yml
name: Deploy API

on:
  push:
    branches: [main]
    paths: ['funkey-bite-api/**']

jobs:
  deploy:
    runs-on: ubuntu-latest
    permissions:
      id-token: write
      contents: read

    steps:
      - uses: actions/checkout@v4

      - name: Configure AWS credentials (OIDC)
        uses: aws-actions/configure-aws-credentials@v4
        with:
          role-to-assume: arn:aws:iam::${{ secrets.AWS_ACCOUNT_ID }}:role/github-actions-deploy
          aws-region: us-east-1

      - name: Login to ECR
        id: login-ecr
        uses: aws-actions/amazon-ecr-login@v2

      - name: Build & push image
        env:
          ECR_URL: ${{ steps.login-ecr.outputs.registry }}/funkey-production-api
          IMAGE_TAG: ${{ github.sha }}
        run: |
          docker build -t $ECR_URL:$IMAGE_TAG funkey-bite-api/
          docker push $ECR_URL:$IMAGE_TAG

      - name: Update ECS service
        run: |
          # Update task definition image tag, then force new deployment
          TASK_DEF=$(aws ecs describe-task-definition \
            --task-definition funkey-production-api \
            --query 'taskDefinition' --output json \
            | jq --arg IMG "$ECR_URL:$IMAGE_TAG" \
              '.containerDefinitions[0].image = $IMG | del(.taskDefinitionArn,.revision,.status,.requiresAttributes,.compatibilities,.registeredAt,.registeredBy)')
          NEW_ARN=$(aws ecs register-task-definition \
            --cli-input-json "$TASK_DEF" \
            --query 'taskDefinition.taskDefinitionArn' --output text)
          aws ecs update-service \
            --cluster funkey-production-api-cluster \
            --service funkey-production-api-service \
            --task-definition "$NEW_ARN" \
            --force-new-deployment
```

> Note: ECS services have `lifecycle { ignore_changes = [task_definition] }` so Terraform won't revert CI/CD image updates.

### GitHub Actions — Frontend Deploy

```yaml
# .github/workflows/deploy-web.yml
name: Deploy Web

on:
  push:
    branches: [main]
    paths: ['funkey-bite-web/**']

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: aws-actions/configure-aws-credentials@v4
        with:
          role-to-assume: arn:aws:iam::${{ secrets.AWS_ACCOUNT_ID }}:role/github-actions-deploy
          aws-region: us-east-1

      - name: Build Nuxt app
        working-directory: funkey-bite-web
        run: npm ci && npm run generate

      - name: Sync to S3 and invalidate CloudFront
        run: |
          aws s3 sync funkey-bite-web/.output/public/ \
            s3://${{ secrets.WEB_BUCKET_ID }}/ --delete
          aws cloudfront create-invalidation \
            --distribution-id ${{ secrets.WEB_CF_DISTRIBUTION_ID }} \
            --paths "/*"
```

---

## Makefile Reference

```
make <target> [ENV=production|staging]

  bootstrap-state         Create S3 state bucket + DynamoDB lock table (run ONCE)

  init                    terraform init -upgrade for ENV
  init-staging            init for staging
  init-production         init for production

  validate                Validate Terraform configuration
  fmt                     Recursively format all .tf files
  fmt-check               Check formatting without changes (CI)

  plan                    Plan ALL components
  plan-infrastructure     Plan networking, WAF, DNS
  plan-database           Plan Aurora + Secrets
  plan-api                Plan ECR + ECS API + ALB
  plan-web                Plan Nuxt web app (S3 + CloudFront)
  plan-admin              Plan React admin (S3 + CloudFront)
  plan-workers            Plan ECS Workers
  plan-lambda             Plan Lambda catering notifier
  plan-cache              Plan ElastiCache Redis
  plan-queue              Plan SQS queues
  plan-scheduler          Plan EventBridge rules
  plan-monitoring         Plan CloudWatch + SNS

  deploy                  Deploy ALL components
  deploy-infrastructure   Deploy networking, WAF, DNS
  deploy-database         Deploy Aurora + Secrets
  deploy-api              Deploy ECR + ECS API + ALB
  deploy-web              Deploy Nuxt web app
  deploy-admin            Deploy React admin dashboard
  deploy-workers          Deploy ECS Workers
  deploy-lambda           Deploy Lambda catering notifier
  deploy-cache            Deploy ElastiCache Redis
  deploy-queue            Deploy SQS queues
  deploy-scheduler        Deploy EventBridge rules
  deploy-monitoring       Deploy CloudWatch + SNS

  deploy-fullstack        infrastructure + database + api
  deploy-frontend-only    web + admin
  deploy-data-plane       database + cache + queue

  output                  Print all outputs for ENV
  output-json             Print outputs as JSON

  state-list              List resources in state
  state-show RESOURCE=X   Show a specific resource
  refresh                 Reconcile state with real resources

  destroy                 Destroy ALL (requires typing ENV name)
  destroy-staging         Destroy all staging resources (-auto-approve)

Default ENV = production. Override: make deploy ENV=staging
```

---

## Security Controls

| Layer | Control |
|---|---|
| **Network** | 3-tier subnets (public/private/data), SGs with least-privilege rules, VPC endpoints for S3/ECR/Secrets (no internet traversal for internal traffic) |
| **Edge** | WAF v2 managed rules (CRS, KnownBadInputs, SQLi), rate-limiting 2 000 req/5 min per IP |
| **ALB** | HTTP→HTTPS redirect; CloudFront secret header validation (direct ALB access returns 403) |
| **CloudFront** | OAC for S3 origins (no public S3 URLs); Origin Shield for cache efficiency and origin protection |
| **Secrets** | All credentials in Secrets Manager; injected into containers at runtime — never baked into images or env var literals; `ignore_changes` prevents Terraform drift |
| **Encryption** | S3: AES-256; Aurora: storage encrypted; ElastiCache: in-transit + at-rest; SQS: SSE-SQS; SNS: KMS-CMK |
| **IAM** | Task-level IAM roles with least-privilege; separate execution vs task roles; no `*` actions except scoped where necessary |
| **Container** | ECR `scan_on_push = true`; lifecycle policies remove old untagged images |
| **Data** | Aurora deletion protection in production; 30-day secret recovery window in production |

---

## Troubleshooting

### ACM certificate stuck in `PENDING_VALIDATION`

Route 53 validation records are created automatically. This usually means NS propagation hasn't completed yet. Check:

```bash
dig NS funkeygrabandbite.com
```

If the nameservers don't match the Route 53 values, update them at your registrar.

### `Error: reference to undeclared resource` when feature flag is false

Cross-module references use `try()` in `locals.tf`. If you see this error, ensure you're referencing `local.X` (from `locals.tf`) rather than `module.X[0].output` directly in `main.tf`.

### ECS service stuck in `PENDING`

1. Check the ECS service events in the AWS Console for the cluster.
2. Common cause: ECR image not pushed yet. Push the image first, then force a new deployment:
   ```bash
   aws ecs update-service --cluster <cluster> --service <service> --force-new-deployment
   ```

### Aurora cluster can't be deleted

Production clusters have `deletion_protection = true`. To destroy:

```bash
# Temporarily disable deletion protection
aws rds modify-db-cluster --db-cluster-identifier funkey-production-aurora \
  --no-deletion-protection --apply-immediately
# Then destroy
make destroy ENV=production
```

### Secrets Manager secret can't be immediately deleted in staging

`recovery_window_in_days = 0` is set for staging but AWS enforces a minimum. If you see an error, use the AWS CLI:

```bash
aws secretsmanager delete-secret \
  --secret-id funkey-staging-app-secrets \
  --force-delete-without-recovery
```

### CloudFront returning 403 for API calls

Verify the secret header is being forwarded correctly:

```bash
# The secret header name and value are sensitive outputs
terraform -chdir=environments/production output -raw cloudfront_origin_secret_header_name
terraform -chdir=environments/production output -raw cloudfront_origin_secret_header_value
```

Confirm the ALB listener rule matches the same header and value.

---

*Built with Terraform · Hosted on AWS · Part of the Funkey Grab & Bite portfolio.*
