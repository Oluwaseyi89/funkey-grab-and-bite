output "web_acl_arn" {
  description = "WAF WebACL ARN (pass to CloudFront distributions)"
  value       = aws_wafv2_web_acl.main.arn
}

output "web_acl_id" {
  description = "WAF WebACL ID"
  value       = aws_wafv2_web_acl.main.id
}
