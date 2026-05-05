locals {
  env         = "staging"
  name_prefix = "${var.project_name}-${local.env}"

  common_tags = {
    Project     = var.project_name
    Environment = local.env
    ManagedBy   = "terraform"
    Owner       = "funkey-team"
  }

  vpc_id             = try(module.networking[0].vpc_id, "")
  public_subnet_ids  = try(module.networking[0].public_subnet_ids, [])
  private_subnet_ids = try(module.networking[0].private_subnet_ids, [])
  data_subnet_ids    = try(module.networking[0].data_subnet_ids, [])
  alb_sg_id          = try(module.networking[0].alb_sg_id, "")
  ecs_api_sg_id      = try(module.networking[0].ecs_api_sg_id, "")
  ecs_workers_sg_id  = try(module.networking[0].ecs_workers_sg_id, "")
  aurora_sg_id       = try(module.networking[0].aurora_sg_id, "")
  elasticache_sg_id  = try(module.networking[0].elasticache_sg_id, "")
  lambda_sg_id       = try(module.networking[0].lambda_sg_id, "")

  certificate_arn = try(module.dns[0].certificate_arn, "")
  zone_id         = try(module.dns[0].zone_id, "")
  waf_web_acl_arn = try(module.waf[0].web_acl_arn, "")
  alb_dns_name    = try(module.alb[0].alb_dns_name, "")

  cf_origin_secret_name  = try(module.cloudfront_api[0].origin_secret_header_name, "X-CloudFront-Secret")
  cf_origin_secret_value = try(module.cloudfront_api[0].origin_secret_header_value, "placeholder")

  db_secret_arn   = try(module.aurora[0].db_secret_arn, "")
  app_secrets_arn = try(module.secrets[0].app_secrets_arn, "")

  ecr_api_url         = try(module.ecr[0].api_repository_url, "")
  lambda_catering_ecr = try(module.ecr[0].lambda_catering_repository_url, "")

  sqs_order_queue_url    = try(module.sqs[0].order_queue_url, "")
  sqs_catering_queue_url = try(module.sqs[0].catering_queue_url, "")
  sqs_catering_queue_arn = try(module.sqs[0].catering_queue_arn, "")
}
