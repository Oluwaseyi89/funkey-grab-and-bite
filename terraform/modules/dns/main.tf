# ─────────────────────────────────────────────────────────────────
# Route 53 Hosted Zone
# ─────────────────────────────────────────────────────────────────

resource "aws_route53_zone" "main" {
  count = var.create_zone ? 1 : 0
  name  = var.domain_name
  tags  = var.tags
}

data "aws_route53_zone" "existing" {
  count   = var.create_zone ? 0 : 1
  zone_id = var.existing_zone_id
}

locals {
  zone_id   = var.create_zone ? aws_route53_zone.main[0].zone_id : data.aws_route53_zone.existing[0].zone_id
  zone_name = var.create_zone ? aws_route53_zone.main[0].name : data.aws_route53_zone.existing[0].name
}

# ─────────────────────────────────────────────────────────────────
# ACM Certificate — wildcard covers all subdomains
# Must be in us-east-1 for CloudFront usage (already enforced at
# provider level in the environment root config)
# ─────────────────────────────────────────────────────────────────

resource "aws_acm_certificate" "main" {
  domain_name               = var.domain_name
  subject_alternative_names = ["*.${var.domain_name}"]
  validation_method         = "DNS"

  tags = merge(var.tags, { Name = "${var.domain_name}-cert" })

  lifecycle {
    create_before_destroy = true
  }
}

# Create DNS validation records in Route 53
resource "aws_route53_record" "cert_validation" {
  for_each = {
    for dvo in aws_acm_certificate.main.domain_validation_options :
    dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  allow_overwrite = true
  name            = each.value.name
  records         = [each.value.record]
  ttl             = 60
  type            = each.value.type
  zone_id         = local.zone_id
}

# Wait for certificate to be validated (can take up to 30 min first time)
resource "aws_acm_certificate_validation" "main" {
  certificate_arn         = aws_acm_certificate.main.arn
  validation_record_fqdns = [for r in aws_route53_record.cert_validation : r.fqdn]
}

# ─────────────────────────────────────────────────────────────────
# Health Check for API endpoint (used by Route 53 failover)
# ─────────────────────────────────────────────────────────────────

resource "aws_route53_health_check" "api" {
  fqdn              = "api.${var.domain_name}"
  port              = 443
  type              = "HTTPS"
  resource_path     = "/api/v1/settings"
  failure_threshold = "3"
  request_interval  = "30"

  tags = merge(var.tags, { Name = "${var.domain_name}-api-health-check" })
}
