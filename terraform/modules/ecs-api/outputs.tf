output "cluster_id" {
  description = "ECS API cluster ID"
  value       = aws_ecs_cluster.api.id
}

output "cluster_arn" {
  description = "ECS API cluster ARN"
  value       = aws_ecs_cluster.api.arn
}

output "cluster_name" {
  description = "ECS API cluster name"
  value       = aws_ecs_cluster.api.name
}

output "service_name" {
  description = "ECS API service name"
  value       = aws_ecs_service.api.name
}

output "task_definition_arn" {
  description = "Task definition ARN (latest revision)"
  value       = aws_ecs_task_definition.api.arn
}

output "execution_role_arn" {
  description = "Task execution IAM role ARN"
  value       = aws_iam_role.ecs_execution.arn
}

output "task_role_arn" {
  description = "Task IAM role ARN"
  value       = aws_iam_role.ecs_task.arn
}

output "log_group_name" {
  description = "CloudWatch log group name"
  value       = aws_cloudwatch_log_group.api.name
}
