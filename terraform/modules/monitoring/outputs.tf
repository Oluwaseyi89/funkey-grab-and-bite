output "sns_topic_arn" {
  description = "SNS alert topic ARN"
  value       = aws_sns_topic.alerts.arn
}

output "dashboard_name" {
  description = "CloudWatch dashboard name"
  value       = aws_cloudwatch_dashboard.main.dashboard_name
}

output "log_group_api" {
  description = "CloudWatch log group for ECS API"
  value       = aws_cloudwatch_log_group.api.name
}

output "log_group_workers" {
  description = "CloudWatch log group for ECS workers"
  value       = aws_cloudwatch_log_group.workers.name
}

output "log_group_lambda_catering" {
  description = "CloudWatch log group for Lambda catering-notify"
  value       = aws_cloudwatch_log_group.lambda_catering.name
}
