output "zone_id" {
  description = "Route 53 hosted zone ID"
  value       = local.zone_id
}

output "zone_name" {
  description = "Hosted zone domain name"
  value       = local.zone_name
}

output "name_servers" {
  description = "Route 53 name servers (update your domain registrar)"
  value       = var.create_zone ? aws_route53_zone.main[0].name_servers : []
}

output "certificate_arn" {
  description = "Validated ACM certificate ARN (for CloudFront + ALB)"
  value       = aws_acm_certificate_validation.main.certificate_arn
}

output "health_check_id" {
  description = "Route 53 health check ID for the API"
  value       = aws_route53_health_check.api.id
}
