variable "name_prefix" {
  description = "Resource name prefix"
  type        = string
}

variable "environment" {
  description = "Deployment environment"
  type        = string
}

variable "alb_dns_name" {
  description = "ALB DNS name (CloudFront custom origin)"
  type        = string
}

variable "acm_certificate_arn" {
  description = "ACM certificate ARN (must be in us-east-1)"
  type        = string
}

variable "domain_aliases" {
  description = "Custom domain aliases (e.g. ['api.funkeygrabandbite.com'])"
  type        = list(string)
  default     = []
}

variable "waf_web_acl_arn" {
  description = "WAF WebACL ARN"
  type        = string
  default     = ""
}

variable "origin_shield_region" {
  description = "Enable CloudFront Origin Shield in this region to reduce load on origin"
  type        = string
  default     = "us-east-1"
}

variable "price_class" {
  description = "CloudFront price class"
  type        = string
  default     = "PriceClass_100"
}

variable "origin_secret_header_value" {
  description = "Secret value injected as a custom origin header to prevent direct ALB access"
  type        = string
  sensitive   = true
}

variable "tags" {
  description = "Tags applied to all resources"
  type        = map(string)
  default     = {}
}
