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

variable "alert_email" {
  description = "Email address to receive operational alerts"
  type        = string
}

variable "ecs_api_cluster_name" {
  description = "ECS API cluster name (for alarms)"
  type        = string
  default     = ""
}

variable "ecs_api_service_name" {
  description = "ECS API service name"
  type        = string
  default     = ""
}

variable "alb_arn_suffix" {
  description = "ALB ARN suffix (for CloudWatch metrics)"
  type        = string
  default     = ""
}

variable "target_group_arn_suffix" {
  description = "Target group ARN suffix"
  type        = string
  default     = ""
}

variable "aurora_cluster_identifier" {
  description = "Aurora cluster identifier (for alarms)"
  type        = string
  default     = ""
}

variable "elasticache_replication_group_id" {
  description = "ElastiCache replication group ID"
  type        = string
  default     = ""
}

variable "tags" {
  description = "Tags applied to all resources"
  type        = map(string)
  default     = {}
}
