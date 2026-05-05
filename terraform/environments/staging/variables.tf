variable "aws_region" {
  type    = string
  default = "us-east-1"
}

variable "project_name" {
  type    = string
  default = "funkey"
}

variable "domain_name" {
  type    = string
  default = "staging.funkeygrabandbite.com"
}

variable "alert_email" {
  type = string
}

variable "api_image_tag" {
  type    = string
  default = "latest"
}

variable "ses_sender_email" {
  type    = string
  default = "noreply@funkeygrabandbite.com"
}

variable "features" {
  description = "Feature flags for staging — defaults only deploy what is needed to showcase"
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
  # Staging default: lean deployment without workers, cache, and scheduler
  default = {
    networking      = true
    waf             = false # WAF adds cost; enable when demoing security
    dns             = true
    web_app         = true
    admin_app       = true
    api             = true
    workers         = false # Not needed for showcase; enable for full demo
    lambda_catering = false # Enable when Lambda func image is ready
    database        = true
    cache           = false # App runs fine without Redis initially
    queue           = true
    scheduler       = false
    secrets         = true
    monitoring      = true
  }
}
