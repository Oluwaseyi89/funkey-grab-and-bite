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
  description = "Data-tier subnet IDs for the ElastiCache subnet group"
  type        = list(string)
}

variable "elasticache_sg_id" {
  description = "Security group ID for ElastiCache"
  type        = string
}

variable "node_type" {
  description = "ElastiCache node type"
  type        = string
  default     = "cache.t4g.micro"
}

variable "num_cache_clusters" {
  description = "Number of cache nodes in the replication group (>= 2 for Multi-AZ)"
  type        = number
  default     = 2
}

variable "engine_version" {
  description = "Redis engine version"
  type        = string
  default     = "7.1"
}

variable "tags" {
  description = "Tags applied to all resources"
  type        = map(string)
  default     = {}
}
