# ── Web App ───────────────────────────────────────────────────────
output "web_url" {
  description = "Public URL for the customer web app"
  value       = var.features.web_app ? "https://${var.domain_name}" : "not deployed"
}

output "web_cloudfront_id" {
  description = "CloudFront distribution ID (needed for cache invalidation on deploy)"
  value       = try(module.cloudfront_web[0].distribution_id, "")
}

output "web_s3_bucket" {
  description = "S3 bucket name for web app assets"
  value       = try(module.static_hosting_web[0].bucket_id, "")
}

# ── Admin App ─────────────────────────────────────────────────────
output "admin_url" {
  description = "Public URL for the admin dashboard"
  value       = var.features.admin_app ? "https://admin.${var.domain_name}" : "not deployed"
}

output "admin_cloudfront_id" {
  description = "CloudFront distribution ID for admin app"
  value       = try(module.cloudfront_admin[0].distribution_id, "")
}

output "admin_s3_bucket" {
  description = "S3 bucket name for admin app assets"
  value       = try(module.static_hosting_admin[0].bucket_id, "")
}

# ── API ───────────────────────────────────────────────────────────
output "api_url" {
  description = "Public URL for the Go API"
  value       = var.features.api ? "https://api.${var.domain_name}" : "not deployed"
}

output "api_cloudfront_id" {
  description = "CloudFront distribution ID for API edge"
  value       = try(module.cloudfront_api[0].distribution_id, "")
}

output "ecr_api_repository_url" {
  description = "ECR URL to push the API image: docker push <url>:<tag>"
  value       = try(module.ecr[0].api_repository_url, "")
}

output "ecr_lambda_catering_repository_url" {
  description = "ECR URL for the Lambda catering-notify image"
  value       = try(module.ecr[0].lambda_catering_repository_url, "")
}

# ── Database ──────────────────────────────────────────────────────
output "aurora_endpoint" {
  description = "Aurora writer endpoint"
  value       = try(module.aurora[0].cluster_endpoint, "")
}

output "aurora_reader_endpoint" {
  description = "Aurora reader endpoint (for reports)"
  value       = try(module.aurora[0].reader_endpoint, "")
}

output "db_secret_arn" {
  description = "Secrets Manager ARN for DB credentials"
  value       = try(module.aurora[0].db_secret_arn, "")
  sensitive   = true
}

# ── Cache ─────────────────────────────────────────────────────────
output "redis_primary_endpoint" {
  description = "Redis primary endpoint"
  value       = try(module.elasticache[0].primary_endpoint_address, "")
}

# ── Queues ────────────────────────────────────────────────────────
output "sqs_order_queue_url" {
  description = "SQS order queue URL"
  value       = try(module.sqs[0].order_queue_url, "")
}

output "sqs_catering_queue_url" {
  description = "SQS catering queue URL"
  value       = try(module.sqs[0].catering_queue_url, "")
}

# ── Monitoring ────────────────────────────────────────────────────
output "sns_alerts_arn" {
  description = "SNS topic ARN for operational alerts"
  value       = try(module.monitoring[0].sns_topic_arn, "")
}

# ── Network ───────────────────────────────────────────────────────
output "vpc_id" {
  description = "VPC ID"
  value       = local.vpc_id
}

output "name_servers" {
  description = "Route 53 name servers — update these at your domain registrar"
  value       = try(module.dns[0].name_servers, [])
}
