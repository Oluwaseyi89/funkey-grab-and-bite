variable "name_prefix" {
  description = "Resource name prefix"
  type        = string
}

variable "environment" {
  description = "Deployment environment"
  type        = string
}

variable "message_retention_seconds" {
  description = "SQS message retention period in seconds"
  type        = number
  default     = 345600 # 4 days
}

variable "visibility_timeout_seconds" {
  description = "SQS visibility timeout (should exceed Lambda/worker processing time)"
  type        = number
  default     = 60
}

variable "tags" {
  description = "Tags applied to all resources"
  type        = map(string)
  default     = {}
}
