variable "name_prefix" {
  description = "Resource name prefix"
  type        = string
}

variable "environment" {
  description = "Deployment environment"
  type        = string
}

variable "rate_limit_threshold" {
  description = "Max requests per 5 minutes per IP before rate-limiting"
  type        = number
  default     = 2000
}

variable "tags" {
  description = "Tags applied to all resources"
  type        = map(string)
  default     = {}
}
