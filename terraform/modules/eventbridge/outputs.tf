output "event_bus_arn" {
  description = "Custom EventBridge bus ARN"
  value       = aws_cloudwatch_event_bus.main.arn
}

output "event_bus_name" {
  description = "Custom EventBridge bus name"
  value       = aws_cloudwatch_event_bus.main.name
}

output "scheduler_role_arn" {
  description = "IAM role ARN for EventBridge Scheduler"
  value       = aws_iam_role.scheduler.arn
}
