variable "name_prefix" {
  description = "Resource name prefix"
  type        = string
}

variable "environment" {
  description = "Deployment environment"
  type        = string
}

variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "vpc_id" {
  description = "VPC ID"
  type        = string
}

variable "private_subnet_ids" {
  description = "Private subnet IDs for worker tasks"
  type        = list(string)
}

variable "ecs_workers_sg_id" {
  description = "Security group ID for worker tasks"
  type        = string
}

variable "ecr_repository_url" {
  description = "ECR repository URL for the worker image (can reuse API image)"
  type        = string
}

variable "image_tag" {
  description = "Docker image tag"
  type        = string
  default     = "latest"
}

variable "task_cpu" {
  description = "ECS task CPU units"
  type        = number
  default     = 256
}

variable "task_memory" {
  description = "ECS task memory in MiB"
  type        = number
  default     = 512
}

variable "desired_count" {
  description = "Number of worker tasks to run"
  type        = number
  default     = 1
}

variable "db_secret_arn" {
  description = "Secrets Manager ARN for DB credentials"
  type        = string
}

variable "app_secrets_arn" {
  description = "Secrets Manager ARN for app-level secrets"
  type        = string
}

variable "sqs_order_queue_url" {
  description = "SQS order queue URL (passed to worker as env var)"
  type        = string
  default     = ""
}

variable "tags" {
  description = "Tags applied to all resources"
  type        = map(string)
  default     = {}
}
