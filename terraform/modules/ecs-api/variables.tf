variable "name_prefix" {
  description = "Resource name prefix"
  type        = string
}

variable "environment" {
  description = "Deployment environment"
  type        = string
}

variable "aws_region" {
  description = "AWS region for CloudWatch log groups"
  type        = string
  default     = "us-east-1"
}

variable "vpc_id" {
  description = "VPC ID"
  type        = string
}

variable "private_subnet_ids" {
  description = "Private subnet IDs for ECS tasks"
  type        = list(string)
}

variable "ecs_api_sg_id" {
  description = "Security group ID for ECS API tasks"
  type        = string
}

variable "target_group_arn" {
  description = "ALB target group ARN to register tasks"
  type        = string
}

variable "ecr_repository_url" {
  description = "ECR repository URL for the API image"
  type        = string
}

variable "image_tag" {
  description = "Docker image tag to deploy"
  type        = string
  default     = "latest"
}

variable "task_cpu" {
  description = "ECS task CPU units (256, 512, 1024, 2048, 4096)"
  type        = number
  default     = 512
}

variable "task_memory" {
  description = "ECS task memory in MiB"
  type        = number
  default     = 1024
}

variable "desired_count" {
  description = "Initial number of running task replicas"
  type        = number
  default     = 1
}

variable "min_capacity" {
  description = "Minimum number of tasks for auto-scaling"
  type        = number
  default     = 1
}

variable "max_capacity" {
  description = "Maximum number of tasks for auto-scaling"
  type        = number
  default     = 4
}

variable "db_secret_arn" {
  description = "Secrets Manager ARN for Aurora DB credentials"
  type        = string
}

variable "app_secrets_arn" {
  description = "Secrets Manager ARN for app-level secrets (JWT, SES, Twilio)"
  type        = string
}

variable "ses_sender_email" {
  description = "SES verified sender email address"
  type        = string
  default     = "noreply@funkeygrabandbite.com"
}

variable "tags" {
  description = "Tags applied to all resources"
  type        = map(string)
  default     = {}
}
