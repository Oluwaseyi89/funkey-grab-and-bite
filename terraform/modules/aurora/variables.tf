variable "name_prefix" {
  description = "Resource name prefix"
  type        = string
}

variable "environment" {
  description = "Deployment environment"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID"
  type        = string
}

variable "data_subnet_ids" {
  description = "Data-tier subnet IDs for the Aurora subnet group"
  type        = list(string)
}

variable "aurora_sg_id" {
  description = "Security group ID for Aurora"
  type        = string
}

variable "engine_version" {
  description = "Aurora PostgreSQL engine version"
  type        = string
  default     = "16.4"
}

variable "database_name" {
  description = "Initial database name"
  type        = string
  default     = "funkey_grab_bite"
}

variable "master_username" {
  description = "Aurora master username"
  type        = string
  default     = "funkey_admin"
}

variable "min_capacity" {
  description = "Aurora Serverless v2 minimum ACU (0.5 = minimum)"
  type        = number
  default     = 0.5
}

variable "max_capacity" {
  description = "Aurora Serverless v2 maximum ACU"
  type        = number
  default     = 4
}

variable "backup_retention_days" {
  description = "Days to retain automated backups"
  type        = number
  default     = 7
}

variable "enable_reader" {
  description = "Add a reader instance for reporting queries (Multi-AZ)"
  type        = bool
  default     = true
}

variable "tags" {
  description = "Tags applied to all resources"
  type        = map(string)
  default     = {}
}
