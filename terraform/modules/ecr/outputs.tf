output "api_repository_url" {
  description = "ECR repository URL for the Go API image"
  value       = aws_ecr_repository.api.repository_url
}

output "api_repository_arn" {
  description = "ECR repository ARN for the Go API image"
  value       = aws_ecr_repository.api.arn
}

output "lambda_catering_repository_url" {
  description = "ECR repository URL for the Lambda catering-notify image"
  value       = aws_ecr_repository.lambda_catering.repository_url
}

output "lambda_catering_repository_arn" {
  description = "ECR repository ARN for the Lambda catering-notify image"
  value       = aws_ecr_repository.lambda_catering.arn
}

output "registry_id" {
  description = "ECR registry (AWS account) ID"
  value       = data.aws_caller_identity.current.account_id
}
