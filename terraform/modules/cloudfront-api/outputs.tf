output "distribution_id" {
  description = "CloudFront distribution ID"
  value       = aws_cloudfront_distribution.main.id
}

output "distribution_arn" {
  description = "CloudFront distribution ARN"
  value       = aws_cloudfront_distribution.main.arn
}

output "cloudfront_domain_name" {
  description = "CloudFront-assigned domain name"
  value       = aws_cloudfront_distribution.main.domain_name
}

output "cloudfront_hosted_zone_id" {
  description = "CloudFront hosted zone ID (for Route 53 ALIAS records)"
  value       = aws_cloudfront_distribution.main.hosted_zone_id
}

output "origin_secret_header_name" {
  description = "Custom header name that the ALB listener rule should check"
  value       = "X-CloudFront-Secret"
}

output "origin_secret_header_value" {
  description = "Secret header value expected by ALB (store in Secrets Manager and rotate)"
  value       = random_password.cf_origin_secret.result
  sensitive   = true
}
