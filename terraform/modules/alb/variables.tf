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

variable "public_subnet_ids" {
  description = "Public subnet IDs for the ALB"
  type        = list(string)
}

variable "alb_sg_id" {
  description = "Security group ID for the ALB"
  type        = string
}

variable "acm_certificate_arn" {
  description = "ACM certificate ARN for HTTPS listener"
  type        = string
}

variable "cloudfront_secret_header_name" {
  description = "Header name that CloudFront sends to verify origin requests"
  type        = string
  default     = "X-CloudFront-Secret"
}

variable "cloudfront_secret_header_value" {
  description = "Expected header value — ALB will return 403 if header is absent"
  type        = string
  sensitive   = true
}

variable "health_check_path" {
  description = "Path used for ALB target group health checks"
  type        = string
  default     = "/api/v1/settings"
}

variable "deregistration_delay" {
  description = "Seconds to drain connections before deregistering targets"
  type        = number
  default     = 30
}

variable "tags" {
  description = "Tags applied to all resources"
  type        = map(string)
  default     = {}
}
