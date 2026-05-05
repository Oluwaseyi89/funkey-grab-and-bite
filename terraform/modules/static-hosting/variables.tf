variable "name_prefix" {
  description = "Resource name prefix"
  type        = string
}

variable "bucket_purpose" {
  description = "Short identifier appended to the bucket name (e.g. web, admin)"
  type        = string
}

variable "environment" {
  description = "Deployment environment"
  type        = string
}

variable "force_destroy" {
  description = "Allow Terraform to destroy the bucket even if it contains objects"
  type        = bool
  default     = false
}

variable "tags" {
  description = "Tags applied to all resources"
  type        = map(string)
  default     = {}
}
