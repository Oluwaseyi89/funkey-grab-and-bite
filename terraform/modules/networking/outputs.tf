output "vpc_id" {
  description = "VPC ID"
  value       = aws_vpc.main.id
}

output "vpc_cidr" {
  description = "VPC CIDR block"
  value       = aws_vpc.main.cidr_block
}

output "public_subnet_ids" {
  description = "Public subnet IDs"
  value       = aws_subnet.public[*].id
}

output "private_subnet_ids" {
  description = "Private subnet IDs (for ECS, Lambda)"
  value       = aws_subnet.private[*].id
}

output "data_subnet_ids" {
  description = "Data subnet IDs (for Aurora, ElastiCache)"
  value       = aws_subnet.data[*].id
}

output "alb_sg_id" {
  description = "ALB security group ID"
  value       = aws_security_group.alb.id
}

output "ecs_api_sg_id" {
  description = "ECS API tasks security group ID"
  value       = aws_security_group.ecs_api.id
}

output "ecs_workers_sg_id" {
  description = "ECS worker tasks security group ID"
  value       = aws_security_group.ecs_workers.id
}

output "aurora_sg_id" {
  description = "Aurora cluster security group ID"
  value       = aws_security_group.aurora.id
}

output "elasticache_sg_id" {
  description = "ElastiCache security group ID"
  value       = aws_security_group.elasticache.id
}

output "lambda_sg_id" {
  description = "Lambda functions security group ID"
  value       = aws_security_group.lambda.id
}

output "availability_zones" {
  description = "Availability zones used"
  value       = local.azs
}
