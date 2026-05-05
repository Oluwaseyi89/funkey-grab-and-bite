variable "name_prefix" {
  description = "Resource name prefix"
  type        = string
}

variable "app_name" {
  description = "Short name for this distribution (e.g. web, admin)"
  type        = string
}

variable "environment" {
  description = "Deployment environment"
  type        = string
}

variable "bucket_regional_domain_name" {
  description = "S3 bucket regional domain name (from static-hosting module)"
  type        = string
}

variable "oac_id" {
  description = "CloudFront Origin Access Control ID (from static-hosting module)"
  type        = string
}

variable "acm_certificate_arn" {
  description = "ACM certificate ARN (must be in us-east-1)"
  type        = string
}

variable "domain_aliases" {
  description = "Custom domain aliases for this distribution (e.g. ['funkeygrabandbite.com', 'www.funkeygrabandbite.com'])"
  type        = list(string)
  default     = []
}

variable "waf_web_acl_arn" {
  description = "WAF WebACL ARN to associate with this distribution"
  type        = string
  default     = ""
}

variable "price_class" {
  description = "CloudFront price class"
  type        = string
  default     = "PriceClass_100"
}

variable "tags" {
  description = "Tags applied to all resources"
  type        = map(string)
  default     = {}
}
