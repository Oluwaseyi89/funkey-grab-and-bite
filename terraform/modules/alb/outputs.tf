output "alb_arn" {
  description = "ALB ARN"
  value       = aws_lb.main.arn
}

output "alb_dns_name" {
  description = "ALB DNS name (pass to CloudFront API origin)"
  value       = aws_lb.main.dns_name
}

output "alb_zone_id" {
  description = "ALB canonical hosted zone ID (for Route 53 ALIAS)"
  value       = aws_lb.main.zone_id
}

output "target_group_arn" {
  description = "API target group ARN (pass to ECS service)"
  value       = aws_lb_target_group.api.arn
}

output "https_listener_arn" {
  description = "HTTPS listener ARN"
  value       = aws_lb_listener.https.arn
}
