# ═══════════════════════════════════════════════════════════════════
# Funkey Grab & Bite — Production Infrastructure
#
# Deploy order for first-time setup:
#   1. make deploy-infrastructure  (networking + WAF + DNS + secrets)
#   2. make deploy-database        (Aurora — takes ~5 min)
#   3. make deploy-api             (ECR + ALB + ECS)
#   4. make deploy-web             (S3 + CloudFront web)
#   5. make deploy-admin           (S3 + CloudFront admin)
#   6. make deploy-queue           (SQS + EventBridge)
#   7. make deploy-lambda          (Lambda catering notify)
#   8. make deploy-cache           (ElastiCache Redis)
#   9. make deploy-monitoring      (CloudWatch + SNS)
#
# Partial deploys: terraform apply -target=module.<name>
# ═══════════════════════════════════════════════════════════════════

# ── 1. Networking ─────────────────────────────────────────────────

module "networking" {
  count  = var.features.networking ? 1 : 0
  source = "../../modules/networking"

  name_prefix        = local.name_prefix
  environment        = local.env
  aws_region         = var.aws_region
  vpc_cidr           = "10.0.0.0/16"
  az_count           = 3
  single_nat_gateway = false # production: one NAT per AZ
}

# ── 2. WAF ────────────────────────────────────────────────────────

module "waf" {
  count  = var.features.waf ? 1 : 0
  source = "../../modules/waf"

  name_prefix          = local.name_prefix
  environment          = local.env
  rate_limit_threshold = 2000
}

# ── 3. DNS + ACM Certificate ──────────────────────────────────────

module "dns" {
  count  = var.features.dns ? 1 : 0
  source = "../../modules/dns"

  domain_name  = var.domain_name
  create_zone  = true
  environment  = local.env
}

# ── 4. ECR Repositories ───────────────────────────────────────────
# Deploy ECR early so CI/CD can push images before ECS is created

module "ecr" {
  count  = var.features.api ? 1 : 0
  source = "../../modules/ecr"

  name_prefix         = local.name_prefix
  environment         = local.env
  image_count_to_keep = 10
}

# ── 5. Secrets Manager ────────────────────────────────────────────

module "secrets" {
  count  = var.features.secrets ? 1 : 0
  source = "../../modules/secrets"

  name_prefix      = local.name_prefix
  environment      = local.env
  ses_sender_email = var.ses_sender_email
}

# ── 6. Aurora PostgreSQL Serverless v2 ────────────────────────────

module "aurora" {
  count  = var.features.database ? 1 : 0
  source = "../../modules/aurora"

  name_prefix           = local.name_prefix
  environment           = local.env
  vpc_id                = local.vpc_id
  data_subnet_ids       = local.data_subnet_ids
  aurora_sg_id          = local.aurora_sg_id
  database_name         = "funkey_grab_bite"
  master_username       = "funkey_admin"
  min_capacity          = 0.5
  max_capacity          = 8
  backup_retention_days = 7
  enable_reader         = true

  depends_on = [module.networking]
}

# ── 7. ElastiCache Redis ──────────────────────────────────────────

module "elasticache" {
  count  = var.features.cache ? 1 : 0
  source = "../../modules/elasticache"

  name_prefix        = local.name_prefix
  environment        = local.env
  vpc_id             = local.vpc_id
  data_subnet_ids    = local.data_subnet_ids
  elasticache_sg_id  = local.elasticache_sg_id
  node_type          = "cache.t4g.small"
  num_cache_clusters = 2
  engine_version     = "7.1"

  depends_on = [module.networking]
}

# ── 8. SQS Queues ─────────────────────────────────────────────────

module "sqs" {
  count  = var.features.queue ? 1 : 0
  source = "../../modules/sqs"

  name_prefix = local.name_prefix
  environment = local.env
}

# ── 9. Static Hosting — Web App (Nuxt) ───────────────────────────

module "static_hosting_web" {
  count  = var.features.web_app ? 1 : 0
  source = "../../modules/static-hosting"

  name_prefix    = local.name_prefix
  bucket_purpose = "web"
  environment    = local.env
  force_destroy  = false
}

# ── 10. CloudFront — Web App ──────────────────────────────────────

module "cloudfront_web" {
  count  = var.features.web_app ? 1 : 0
  source = "../../modules/cloudfront-spa"

  name_prefix                 = local.name_prefix
  app_name                    = "web"
  environment                 = local.env
  bucket_regional_domain_name = module.static_hosting_web[0].bucket_regional_domain_name
  oac_id                      = module.static_hosting_web[0].oac_id
  acm_certificate_arn         = local.certificate_arn
  domain_aliases              = [var.domain_name, "www.${var.domain_name}"]
  waf_web_acl_arn             = local.waf_web_acl_arn

  depends_on = [module.dns, module.waf]
}

# ── 11. Static Hosting — Admin App (React) ───────────────────────

module "static_hosting_admin" {
  count  = var.features.admin_app ? 1 : 0
  source = "../../modules/static-hosting"

  name_prefix    = local.name_prefix
  bucket_purpose = "admin"
  environment    = local.env
  force_destroy  = false
}

# ── 12. CloudFront — Admin App ────────────────────────────────────

module "cloudfront_admin" {
  count  = var.features.admin_app ? 1 : 0
  source = "../../modules/cloudfront-spa"

  name_prefix                 = local.name_prefix
  app_name                    = "admin"
  environment                 = local.env
  bucket_regional_domain_name = module.static_hosting_admin[0].bucket_regional_domain_name
  oac_id                      = module.static_hosting_admin[0].oac_id
  acm_certificate_arn         = local.certificate_arn
  domain_aliases              = ["admin.${var.domain_name}"]
  waf_web_acl_arn             = local.waf_web_acl_arn

  depends_on = [module.dns, module.waf]
}

# ── 13. CloudFront — API Edge ─────────────────────────────────────

module "cloudfront_api" {
  count  = var.features.api ? 1 : 0
  source = "../../modules/cloudfront-api"

  name_prefix         = local.name_prefix
  environment         = local.env
  alb_dns_name        = local.alb_dns_name
  acm_certificate_arn = local.certificate_arn
  domain_aliases      = ["api.${var.domain_name}"]
  waf_web_acl_arn     = local.waf_web_acl_arn

  depends_on = [module.alb, module.dns, module.waf]
}

# ── 14. ALB ───────────────────────────────────────────────────────

module "alb" {
  count  = var.features.api ? 1 : 0
  source = "../../modules/alb"

  name_prefix                    = local.name_prefix
  environment                    = local.env
  vpc_id                         = local.vpc_id
  public_subnet_ids              = local.public_subnet_ids
  alb_sg_id                      = local.alb_sg_id
  acm_certificate_arn            = local.certificate_arn
  cloudfront_secret_header_name  = local.cf_origin_secret_name
  cloudfront_secret_header_value = local.cf_origin_secret_value
  health_check_path              = "/api/v1/settings"

  depends_on = [module.networking, module.dns, module.cloudfront_api]
}

# ── 15. ECS — API Cluster ─────────────────────────────────────────

module "ecs_api" {
  count  = var.features.api ? 1 : 0
  source = "../../modules/ecs-api"

  name_prefix        = local.name_prefix
  environment        = local.env
  aws_region         = var.aws_region
  vpc_id             = local.vpc_id
  private_subnet_ids = local.private_subnet_ids
  ecs_api_sg_id      = local.ecs_api_sg_id
  target_group_arn   = module.alb[0].target_group_arn
  ecr_repository_url = local.ecr_api_url
  image_tag          = var.api_image_tag
  task_cpu           = 512
  task_memory        = 1024
  desired_count      = 2
  min_capacity       = 2
  max_capacity       = 8
  db_secret_arn      = local.db_secret_arn
  app_secrets_arn    = local.app_secrets_arn
  ses_sender_email   = var.ses_sender_email

  depends_on = [module.networking, module.aurora, module.secrets, module.alb]
}

# ── 16. ECS — Worker Cluster ──────────────────────────────────────

module "ecs_workers" {
  count  = var.features.workers ? 1 : 0
  source = "../../modules/ecs-workers"

  name_prefix         = local.name_prefix
  environment         = local.env
  aws_region          = var.aws_region
  vpc_id              = local.vpc_id
  private_subnet_ids  = local.private_subnet_ids
  ecs_workers_sg_id   = local.ecs_workers_sg_id
  ecr_repository_url  = local.ecr_api_url
  image_tag           = var.api_image_tag
  task_cpu            = 256
  task_memory         = 512
  desired_count       = 1
  db_secret_arn       = local.db_secret_arn
  app_secrets_arn     = local.app_secrets_arn
  sqs_order_queue_url = local.sqs_order_queue_url

  depends_on = [module.networking, module.aurora, module.secrets, module.sqs]
}

# ── 17. Lambda — Catering Notify ─────────────────────────────────

module "lambda_catering" {
  count  = var.features.lambda_catering ? 1 : 0
  source = "../../modules/lambda-catering"

  name_prefix            = local.name_prefix
  environment            = local.env
  aws_region             = var.aws_region
  vpc_id                 = local.vpc_id
  private_subnet_ids     = local.private_subnet_ids
  lambda_sg_id           = local.lambda_sg_id
  ecr_repository_url     = local.lambda_catering_ecr
  image_tag              = var.api_image_tag
  sqs_catering_queue_arn = local.sqs_catering_queue_arn
  sqs_catering_queue_url = local.sqs_catering_queue_url
  app_secrets_arn        = local.app_secrets_arn
  ses_sender_email       = var.ses_sender_email
  reserved_concurrency   = 5

  depends_on = [module.networking, module.sqs, module.secrets, module.ecr]
}

# ── 18. EventBridge Scheduler ─────────────────────────────────────

module "eventbridge" {
  count  = var.features.scheduler ? 1 : 0
  source = "../../modules/eventbridge"

  name_prefix                 = local.name_prefix
  environment                 = local.env
  ecs_api_cluster_arn         = try(module.ecs_api[0].cluster_arn, "")
  ecs_api_task_definition_arn = try(module.ecs_api[0].task_definition_arn, "")
  ecs_task_role_arn           = try(module.ecs_api[0].task_role_arn, "")
  private_subnet_ids          = local.private_subnet_ids
  ecs_api_sg_id               = local.ecs_api_sg_id

  depends_on = [module.ecs_api]
}

# ── 19. Monitoring ────────────────────────────────────────────────

module "monitoring" {
  count  = var.features.monitoring ? 1 : 0
  source = "../../modules/monitoring"

  name_prefix                      = local.name_prefix
  environment                      = local.env
  aws_region                       = var.aws_region
  alert_email                      = var.alert_email
  ecs_api_cluster_name             = try(module.ecs_api[0].cluster_name, "")
  ecs_api_service_name             = try(module.ecs_api[0].service_name, "")
  alb_arn_suffix                   = try(split(":", module.alb[0].alb_arn)[5], "")
  target_group_arn_suffix          = try(split(":", module.alb[0].target_group_arn)[5], "")
  aurora_cluster_identifier        = try(module.aurora[0].cluster_identifier, "")
  elasticache_replication_group_id = try(module.elasticache[0].replication_group_id, "")
}

# ── Route 53 Records (wires everything together) ──────────────────

resource "aws_route53_record" "web_root" {
  count   = (var.features.dns && var.features.web_app) ? 1 : 0
  zone_id = local.zone_id
  name    = var.domain_name
  type    = "A"

  alias {
    name                   = module.cloudfront_web[0].cloudfront_domain_name
    zone_id                = module.cloudfront_web[0].cloudfront_hosted_zone_id
    evaluate_target_health = false
  }
}

resource "aws_route53_record" "web_www" {
  count   = (var.features.dns && var.features.web_app) ? 1 : 0
  zone_id = local.zone_id
  name    = "www.${var.domain_name}"
  type    = "A"

  alias {
    name                   = module.cloudfront_web[0].cloudfront_domain_name
    zone_id                = module.cloudfront_web[0].cloudfront_hosted_zone_id
    evaluate_target_health = false
  }
}

resource "aws_route53_record" "admin" {
  count   = (var.features.dns && var.features.admin_app) ? 1 : 0
  zone_id = local.zone_id
  name    = "admin.${var.domain_name}"
  type    = "A"

  alias {
    name                   = module.cloudfront_admin[0].cloudfront_domain_name
    zone_id                = module.cloudfront_admin[0].cloudfront_hosted_zone_id
    evaluate_target_health = false
  }
}

resource "aws_route53_record" "api" {
  count   = (var.features.dns && var.features.api) ? 1 : 0
  zone_id = local.zone_id
  name    = "api.${var.domain_name}"
  type    = "A"

  alias {
    name                   = module.cloudfront_api[0].cloudfront_domain_name
    zone_id                = module.cloudfront_api[0].cloudfront_hosted_zone_id
    evaluate_target_health = false
  }
}
