output "web_url" {
  value = var.features.web_app ? "https://${var.domain_name}" : "not deployed"
}
output "admin_url" {
  value = var.features.admin_app ? "https://admin.${var.domain_name}" : "not deployed"
}
output "api_url" {
  value = var.features.api ? "https://api.${var.domain_name}" : "not deployed"
}
output "web_cloudfront_id" {
  value = try(module.cloudfront_web[0].distribution_id, "")
}
output "admin_cloudfront_id" {
  value = try(module.cloudfront_admin[0].distribution_id, "")
}
output "web_s3_bucket" {
  value = try(module.static_hosting_web[0].bucket_id, "")
}
output "admin_s3_bucket" {
  value = try(module.static_hosting_admin[0].bucket_id, "")
}
output "ecr_api_repository_url" {
  value = try(module.ecr[0].api_repository_url, "")
}
output "aurora_endpoint" {
  value = try(module.aurora[0].cluster_endpoint, "")
}
output "vpc_id" {
  value = local.vpc_id
}
output "name_servers" {
  description = "Update these at your domain registrar"
  value       = try(module.dns[0].name_servers, [])
}
