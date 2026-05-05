output "cluster_endpoint" {
  description = "Aurora writer endpoint"
  value       = aws_rds_cluster.main.endpoint
}

output "reader_endpoint" {
  description = "Aurora reader endpoint (for reporting)"
  value       = aws_rds_cluster.main.reader_endpoint
}

output "cluster_identifier" {
  description = "Aurora cluster identifier"
  value       = aws_rds_cluster.main.cluster_identifier
}

output "port" {
  description = "Aurora port"
  value       = aws_rds_cluster.main.port
}

output "database_name" {
  description = "Database name"
  value       = aws_rds_cluster.main.database_name
}

output "db_secret_arn" {
  description = "Secrets Manager ARN containing Aurora credentials (host/port/user/pass/dbname)"
  value       = aws_secretsmanager_secret.db_credentials.arn
  sensitive   = true
}
