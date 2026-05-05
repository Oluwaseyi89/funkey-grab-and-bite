variable "name_prefix" {
  description = "Resource name prefix"
  type        = string
}

variable "environment" {
  description = "Deployment environment"
  type        = string
}

# Application-level secrets (JWT, SES, Twilio)
variable "jwt_secret_initial_value" {
  description = "Initial JWT secret value (rotate immediately after first deploy)"
  type        = string
  sensitive   = true
  default     = "REPLACE_ME_JWT_SECRET_32_CHARS_MIN"
}

variable "ses_sender_email" {
  description = "SES verified sender email address"
  type        = string
  default     = "noreply@funkeygrabandbite.com"
}

variable "tags" {
  description = "Tags applied to all resources"
  type        = map(string)
  default     = {}
}
