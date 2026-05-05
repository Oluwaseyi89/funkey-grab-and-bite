variable "name_prefix" {
  description = "Prefix used for naming all resources (e.g. funkey-prod)"
  type        = string
}

variable "environment" {
  description = "Deployment environment (production, staging)"
  type        = string
}

variable "aws_region" {
  description = "AWS region (used to build VPC endpoint service names)"
  type        = string
  default     = "us-east-1"
}

variable "vpc_cidr" {
  description = "CIDR block for the VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "az_count" {
  description = "Number of availability zones to deploy across (2 or 3)"
  type        = number
  default     = 2

  validation {
    condition     = var.az_count >= 2 && var.az_count <= 3
    error_message = "az_count must be 2 or 3."
  }
}

variable "single_nat_gateway" {
  description = "Use a single shared NAT gateway to reduce cost (recommended for staging)"
  type        = bool
  default     = false
}

variable "tags" {
  description = "Tags applied to all resources"
  type        = map(string)
  default     = {}
}
