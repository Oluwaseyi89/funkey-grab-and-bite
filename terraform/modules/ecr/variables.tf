variable "name_prefix" {
  description = "Resource name prefix"
  type        = string
}

variable "environment" {
  description = "Deployment environment"
  type        = string
}

variable "image_count_to_keep" {
  description = "Number of most-recent image tags to retain in ECR"
  type        = number
  default     = 10
}

variable "tags" {
  description = "Tags applied to all resources"
  type        = map(string)
  default     = {}
}
