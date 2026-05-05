output "function_arn" {
  description = "Lambda function ARN"
  value       = aws_lambda_function.catering_notify.arn
}

output "function_name" {
  description = "Lambda function name"
  value       = aws_lambda_function.catering_notify.function_name
}

output "log_group_name" {
  description = "CloudWatch log group name"
  value       = aws_cloudwatch_log_group.lambda.name
}

output "role_arn" {
  description = "Lambda IAM role ARN"
  value       = aws_iam_role.lambda.arn
}
