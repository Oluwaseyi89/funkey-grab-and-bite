.DEFAULT_GOAL := help
SHELL         := /bin/bash

# ─── Environment ─────────────────────────────────────────────────────────────
ENV     ?= production
ENV_DIR  = environments/$(ENV)

# Validate ENV
ifeq ($(filter $(ENV),production staging),)
$(error ENV must be 'production' or 'staging' — got '$(ENV)')
endif

# ─── Colours ──────────────────────────────────────────────────────────────────
BLUE   := \033[0;34m
GREEN  := \033[0;32m
YELLOW := \033[0;33m
RESET  := \033[0m

# ─── Bootstrap (one-time) ─────────────────────────────────────────────────────
.PHONY: bootstrap-state
bootstrap-state: ## Create S3 state bucket + DynamoDB lock table (run ONCE before init)
	@echo -e "$(BLUE)Bootstrapping remote state for environment: $(ENV)$(RESET)"
	@chmod +x scripts/bootstrap-state.sh && ./scripts/bootstrap-state.sh

# ─── Init ─────────────────────────────────────────────────────────────────────
.PHONY: init
init: ## terraform init -upgrade for ENV (default: production)
	@echo -e "$(BLUE)Initialising Terraform for environment: $(ENV)$(RESET)"
	terraform -chdir=$(ENV_DIR) init -upgrade

.PHONY: init-staging
init-staging: ## terraform init for staging
	@$(MAKE) init ENV=staging

.PHONY: init-production
init-production: ## terraform init for production
	@$(MAKE) init ENV=production

# ─── Validate / Format ────────────────────────────────────────────────────────
.PHONY: validate
validate: ## Validate all Terraform configurations
	@echo -e "$(BLUE)Validating configuration for environment: $(ENV)$(RESET)"
	terraform -chdir=$(ENV_DIR) validate

.PHONY: fmt
fmt: ## Recursively format all .tf files
	@echo -e "$(BLUE)Formatting Terraform files...$(RESET)"
	terraform fmt -recursive

.PHONY: fmt-check
fmt-check: ## Check formatting without making changes (CI)
	terraform fmt -recursive -check

# ─── Plan ─────────────────────────────────────────────────────────────────────
.PHONY: plan
plan: ## Plan ALL components for ENV
	@chmod +x scripts/plan.sh && ./scripts/plan.sh $(ENV) all

.PHONY: plan-infrastructure
plan-infrastructure: ## Plan networking, WAF, DNS
	@chmod +x scripts/plan.sh && ./scripts/plan.sh $(ENV) infrastructure

.PHONY: plan-database
plan-database: ## Plan Aurora + Secrets
	@chmod +x scripts/plan.sh && ./scripts/plan.sh $(ENV) database

.PHONY: plan-api
plan-api: ## Plan ECR + ECS API + ALB
	@chmod +x scripts/plan.sh && ./scripts/plan.sh $(ENV) api

.PHONY: plan-web
plan-web: ## Plan S3 + CloudFront for Nuxt web app
	@chmod +x scripts/plan.sh && ./scripts/plan.sh $(ENV) web

.PHONY: plan-admin
plan-admin: ## Plan S3 + CloudFront for React admin
	@chmod +x scripts/plan.sh && ./scripts/plan.sh $(ENV) admin

.PHONY: plan-workers
plan-workers: ## Plan ECS Workers cluster
	@chmod +x scripts/plan.sh && ./scripts/plan.sh $(ENV) workers

.PHONY: plan-lambda
plan-lambda: ## Plan Lambda catering notifier
	@chmod +x scripts/plan.sh && ./scripts/plan.sh $(ENV) lambda

.PHONY: plan-cache
plan-cache: ## Plan ElastiCache Redis
	@chmod +x scripts/plan.sh && ./scripts/plan.sh $(ENV) cache

.PHONY: plan-queue
plan-queue: ## Plan SQS queues + DLQs
	@chmod +x scripts/plan.sh && ./scripts/plan.sh $(ENV) queue

.PHONY: plan-scheduler
plan-scheduler: ## Plan EventBridge Scheduler rules
	@chmod +x scripts/plan.sh && ./scripts/plan.sh $(ENV) scheduler

.PHONY: plan-monitoring
plan-monitoring: ## Plan CloudWatch alarms + SNS + dashboard
	@chmod +x scripts/plan.sh && ./scripts/plan.sh $(ENV) monitoring

# ─── Deploy ───────────────────────────────────────────────────────────────────
.PHONY: deploy
deploy: ## Deploy ALL components for ENV
	@chmod +x scripts/deploy.sh && ./scripts/deploy.sh $(ENV) all

.PHONY: deploy-infrastructure
deploy-infrastructure: ## Deploy networking, WAF, DNS
	@chmod +x scripts/deploy.sh && ./scripts/deploy.sh $(ENV) infrastructure

.PHONY: deploy-database
deploy-database: ## Deploy Aurora + Secrets
	@chmod +x scripts/deploy.sh && ./scripts/deploy.sh $(ENV) database

.PHONY: deploy-api
deploy-api: ## Deploy ECR + ECS API + ALB
	@chmod +x scripts/deploy.sh && ./scripts/deploy.sh $(ENV) api

.PHONY: deploy-web
deploy-web: ## Deploy Nuxt web app (S3 + CloudFront)
	@chmod +x scripts/deploy.sh && ./scripts/deploy.sh $(ENV) web

.PHONY: deploy-admin
deploy-admin: ## Deploy React admin dashboard (S3 + CloudFront)
	@chmod +x scripts/deploy.sh && ./scripts/deploy.sh $(ENV) admin

.PHONY: deploy-workers
deploy-workers: ## Deploy ECS Worker tasks cluster
	@chmod +x scripts/deploy.sh && ./scripts/deploy.sh $(ENV) workers

.PHONY: deploy-lambda
deploy-lambda: ## Deploy Lambda catering notifier
	@chmod +x scripts/deploy.sh && ./scripts/deploy.sh $(ENV) lambda

.PHONY: deploy-cache
deploy-cache: ## Deploy ElastiCache Redis
	@chmod +x scripts/deploy.sh && ./scripts/deploy.sh $(ENV) cache

.PHONY: deploy-queue
deploy-queue: ## Deploy SQS queues + DLQs
	@chmod +x scripts/deploy.sh && ./scripts/deploy.sh $(ENV) queue

.PHONY: deploy-scheduler
deploy-scheduler: ## Deploy EventBridge Scheduler rules
	@chmod +x scripts/deploy.sh && ./scripts/deploy.sh $(ENV) scheduler

.PHONY: deploy-monitoring
deploy-monitoring: ## Deploy CloudWatch alarms + SNS + dashboard
	@chmod +x scripts/deploy.sh && ./scripts/deploy.sh $(ENV) monitoring

# ─── Outputs ──────────────────────────────────────────────────────────────────
.PHONY: output
output: ## Print all Terraform outputs for ENV
	@echo -e "$(BLUE)Outputs for environment: $(ENV)$(RESET)"
	terraform -chdir=$(ENV_DIR) output

.PHONY: output-json
output-json: ## Print all outputs as JSON (useful for CI/CD)
	terraform -chdir=$(ENV_DIR) output -json

# ─── State Inspection ─────────────────────────────────────────────────────────
.PHONY: state-list
state-list: ## List all resources in state for ENV
	terraform -chdir=$(ENV_DIR) state list

.PHONY: state-show
state-show: ## Show a specific resource (usage: make state-show RESOURCE=module.vpc.aws_vpc.main)
	@if [ -z "$(RESOURCE)" ]; then echo "Usage: make state-show RESOURCE=<address>"; exit 1; fi
	terraform -chdir=$(ENV_DIR) state show $(RESOURCE)

# ─── Targeted Refresh ─────────────────────────────────────────────────────────
.PHONY: refresh
refresh: ## Reconcile state with real-world resources for ENV
	terraform -chdir=$(ENV_DIR) apply -refresh-only -auto-approve

# ─── Destroy ──────────────────────────────────────────────────────────────────
.PHONY: destroy
destroy: ## Destroy ALL resources for ENV (REQUIRES confirmation)
	@echo -e "$(YELLOW)WARNING: This will destroy ALL infrastructure for ENV=$(ENV)$(RESET)"
	@read -p "Type the environment name to confirm: " CONFIRM && \
	  if [ "$$CONFIRM" != "$(ENV)" ]; then echo "Aborted."; exit 1; fi
	terraform -chdir=$(ENV_DIR) destroy

.PHONY: destroy-staging
destroy-staging: ## Destroy ALL staging resources (no confirmation prompt on staging)
	@echo -e "$(YELLOW)Destroying all staging resources...$(RESET)"
	terraform -chdir=environments/staging destroy -auto-approve

# ─── Convenience combos ───────────────────────────────────────────────────────
.PHONY: deploy-fullstack
deploy-fullstack: ## Deploy infrastructure + database + api (core stack)
	@$(MAKE) deploy-infrastructure && \
	 $(MAKE) deploy-database       && \
	 $(MAKE) deploy-api

.PHONY: deploy-frontend-only
deploy-frontend-only: ## Deploy web + admin apps only (no backend changes)
	@$(MAKE) deploy-web && $(MAKE) deploy-admin

.PHONY: deploy-data-plane
deploy-data-plane: ## Deploy database + cache + queue (data-layer only)
	@$(MAKE) deploy-database && \
	 $(MAKE) deploy-cache    && \
	 $(MAKE) deploy-queue

.PHONY: plan-all-staging
plan-all-staging: ## Plan everything in staging
	@$(MAKE) plan ENV=staging

.PHONY: plan-all-production
plan-all-production: ## Plan everything in production
	@$(MAKE) plan ENV=production

# ─── Help ─────────────────────────────────────────────────────────────────────
.PHONY: help
help: ## List all available make targets with descriptions
	@echo ""
	@echo -e "$(BLUE)Funkey Grab & Bite — Terraform Makefile$(RESET)"
	@echo -e "Usage: make <target> [ENV=production|staging]\n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	  | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-30s$(RESET) %s\n", $$1, $$2}'
	@echo ""
	@echo -e "  Default ENV = $(YELLOW)$(ENV)$(RESET)"
	@echo -e "  Override:  make deploy ENV=staging\n"
