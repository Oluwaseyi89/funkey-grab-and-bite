output "cluster_id" {
  value = aws_ecs_cluster.workers.id
}

output "cluster_arn" {
  value = aws_ecs_cluster.workers.arn
}

output "cluster_name" {
  value = aws_ecs_cluster.workers.name
}

output "service_name" {
  value = aws_ecs_service.workers.name
}

output "task_role_arn" {
  description = "Task IAM role ARN (for EventBridge scheduler)"
  value       = aws_iam_role.ecs_task.arn
}

output "log_group_name" {
  value = aws_cloudwatch_log_group.workers.name
}
