variable "aws_region" {
  description = "AWS region for deployment"
  type        = string
  default     = "us-east-1"
}

variable "project_name" {
  description = "Short project identifier used in resource names"
  type        = string
  default     = "funkey"
}

variable "domain_name" {
  description = "Root domain name"
  type        = string
  default     = "funkeygrabandbite.com"
}

variable "alert_email" {
  description = "Email address for CloudWatch / SNS alerts"
  type        = string
}

variable "api_image_tag" {
  description = "ECS/Lambda image tag to deploy (updated by CI/CD)"
  type        = string
  default     = "latest"
}

variable "ses_sender_email" {
  description = "SES verified sender address"
  type        = string
  default     = "noreply@funkeygrabandbite.com"
}

# ── Feature Flags ─────────────────────────────────────────────────
# Set any flag to false to skip deploying that component.
# Use -target=module.<name> for targeted deploys.

variable "features" {
  description = "Feature flags — set individual components to false to skip deployment"
  type = object({
    networking      = bool
    waf             = bool
    dns             = bool
    web_app         = bool
    admin_app       = bool
    api             = bool
    workers         = bool
    lambda_catering = bool
    database        = bool
    cache           = bool
    queue           = bool
    scheduler       = bool
    secrets         = bool
    monitoring      = bool
  })
  default = {
    networking      = true
    waf             = true
    dns             = true
    web_app         = true
    admin_app       = true
    api             = true
    workers         = true
    lambda_catering = true
    database        = true
    cache           = true
    queue           = true
    scheduler       = true
    secrets         = true
    monitoring      = true
  }
}
