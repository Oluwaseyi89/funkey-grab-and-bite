variable "domain_name" {
  description = "Root domain name (e.g. funkeygrabandbite.com)"
  type        = string
}

variable "create_zone" {
  description = "Whether to create a new Route 53 hosted zone (false = use existing)"
  type        = bool
  default     = true
}

variable "existing_zone_id" {
  description = "Existing Route 53 zone ID (used when create_zone = false)"
  type        = string
  default     = ""
}

variable "environment" {
  description = "Deployment environment"
  type        = string
}

variable "tags" {
  description = "Tags applied to all resources"
  type        = map(string)
  default     = {}
}
