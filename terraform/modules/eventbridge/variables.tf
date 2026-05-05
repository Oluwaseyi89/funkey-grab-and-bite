variable "name_prefix" {
  description = "Resource name prefix"
  type        = string
}

variable "environment" {
  description = "Deployment environment"
  type        = string
}

variable "ecs_api_cluster_arn" {
  description = "ECS API cluster ARN (for scheduled tasks)"
  type        = string
  default     = ""
}

variable "ecs_api_task_definition_arn" {
  description = "ECS task definition ARN for scheduled tasks"
  type        = string
  default     = ""
}

variable "ecs_task_role_arn" {
  description = "ECS task IAM role ARN"
  type        = string
  default     = ""
}

variable "private_subnet_ids" {
  description = "Private subnet IDs for scheduled ECS tasks"
  type        = list(string)
  default     = []
}

variable "ecs_api_sg_id" {
  description = "ECS API security group ID for scheduled tasks"
  type        = string
  default     = ""
}

variable "tags" {
  description = "Tags applied to all resources"
  type        = map(string)
  default     = {}
}
