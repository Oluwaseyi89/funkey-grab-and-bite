output "app_secrets_arn" {
  description = "Secrets Manager ARN for the app secrets bundle (JWT, SES, Twilio)"
  value       = aws_secretsmanager_secret.app_secrets.arn
  sensitive   = true
}

output "paystack_secret_arn" {
  description = "Secrets Manager ARN for Paystack credentials"
  value       = aws_secretsmanager_secret.paystack.arn
  sensitive   = true
}

output "admin_bootstrap_secret_arn" {
  description = "Secrets Manager ARN for admin bootstrap credentials"
  value       = aws_secretsmanager_secret.admin_bootstrap.arn
  sensitive   = true
}
