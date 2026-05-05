# ─────────────────────────────────────────────────────────────────
# Application Secrets
# The secret values are placeholders. Populate them via:
#   aws secretsmanager put-secret-value --secret-id <arn> --secret-string '{"key":"value"}'
# or via the AWS Console before the first ECS + Lambda deployment.
# ─────────────────────────────────────────────────────────────────

locals {
  recovery_window = var.environment == "production" ? 30 : 0
}

# ── App secrets bundle (JWT + SES + Twilio) ───────────────────────
# ECS task definition reads individual keys from this secret using
# the JSON key syntax: <ARN>:<key>::

resource "aws_secretsmanager_secret" "app_secrets" {
  name                    = "${var.name_prefix}/app/secrets"
  description             = "JWT secret, SES config, and Twilio credentials for ${var.name_prefix}"
  recovery_window_in_days = local.recovery_window
  tags                    = var.tags
}

resource "aws_secretsmanager_secret_version" "app_secrets" {
  secret_id = aws_secretsmanager_secret.app_secrets.id

  # Initial placeholder — replace values via CLI or Console before deploy
  secret_string = jsonencode({
    JWT_SECRET            = var.jwt_secret_initial_value
    SES_SENDER_EMAIL      = var.ses_sender_email
    TWILIO_ACCOUNT_SID    = "REPLACE_ME"
    TWILIO_AUTH_TOKEN     = "REPLACE_ME"
    TWILIO_PHONE_NUMBER   = "REPLACE_ME"
    DEFAULT_ADMIN_EMAIL   = "admin@funkeygrabandbite.com"
    DEFAULT_ADMIN_USERNAME = "admin"
    DEFAULT_ADMIN_PASSWORD = "REPLACE_ME_STRONG_PASSWORD"
    DEFAULT_ADMIN_ROLE    = "admin"
    AWS_REGION            = "us-east-1"
  })

  # Prevent Terraform from overwriting manually updated values on re-apply
  lifecycle {
    ignore_changes = [secret_string]
  }
}

# ── Paystack API credentials ──────────────────────────────────────

resource "aws_secretsmanager_secret" "paystack" {
  name                    = "${var.name_prefix}/payments/paystack"
  description             = "Paystack API credentials for ${var.name_prefix}"
  recovery_window_in_days = local.recovery_window
  tags                    = var.tags
}

resource "aws_secretsmanager_secret_version" "paystack" {
  secret_id = aws_secretsmanager_secret.paystack.id

  secret_string = jsonencode({
    PAYSTACK_SECRET_KEY = "sk_live_REPLACE_ME"
    PAYSTACK_PUBLIC_KEY = "pk_live_REPLACE_ME"
  })

  lifecycle {
    ignore_changes = [secret_string]
  }
}

# ── Admin bootstrap credentials ───────────────────────────────────
# Separate secret scoped to initial setup; can be deleted after first boot

resource "aws_secretsmanager_secret" "admin_bootstrap" {
  name                    = "${var.name_prefix}/admin/bootstrap"
  description             = "First-run admin bootstrap credentials"
  recovery_window_in_days = local.recovery_window
  tags                    = var.tags
}

resource "aws_secretsmanager_secret_version" "admin_bootstrap" {
  secret_id = aws_secretsmanager_secret.admin_bootstrap.id

  secret_string = jsonencode({
    DEFAULT_ADMIN_EMAIL    = "admin@funkeygrabandbite.com"
    DEFAULT_ADMIN_USERNAME = "admin"
    DEFAULT_ADMIN_PASSWORD = "REPLACE_ME_STRONG_PASSWORD"
    DEFAULT_ADMIN_ROLE     = "admin"
  })

  lifecycle {
    ignore_changes = [secret_string]
  }
}
