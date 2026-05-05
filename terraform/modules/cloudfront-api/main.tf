# CloudFront distribution for the Go API (ALB origin)
# Features:
#   - ALB origin with custom verification header (prevents direct ALB access)
#   - No caching for API responses (CachingDisabled managed policy)
#   - All headers/cookies/query strings forwarded to origin
#   - WebSocket support (HTTP upgrade handled transparently)
#   - WAF ACL attached

locals {
  origin_id                 = "${var.name_prefix}-alb-origin"
  origin_custom_header_name = "X-CloudFront-Secret"
}

resource "aws_cloudfront_distribution" "main" {
  enabled         = true
  is_ipv6_enabled = true
  price_class     = var.price_class
  aliases         = var.domain_aliases
  web_acl_id      = var.waf_web_acl_arn != "" ? var.waf_web_acl_arn : null
  comment         = "Funkey API edge — ${var.environment}"

  origin {
    domain_name = var.alb_dns_name
    origin_id   = local.origin_id

    # Used by ALB listener rule to ensure only CloudFront can reach the ALB
    custom_header {
      name  = local.origin_custom_header_name
      value = var.origin_secret_header_value
    }

    custom_origin_config {
      http_port              = 80
      https_port             = 443
      origin_protocol_policy = "https-only"
      origin_ssl_protocols   = ["TLSv1.2"]
    }

    origin_shield {
      enabled              = true
      origin_shield_region = var.origin_shield_region
    }
  }

  default_cache_behavior {
    target_origin_id       = local.origin_id
    viewer_protocol_policy = "redirect-to-https"
    allowed_methods        = ["DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT"]
    cached_methods         = ["GET", "HEAD"]
    compress               = true

    # AWS managed CachingDisabled policy (pass-through for all API requests)
    cache_policy_id = "4135ea2d-6df8-44a3-9df3-4b5a84be39ad"

    # AWS managed AllViewer origin request policy
    # (forward all headers, cookies, query strings to origin)
    origin_request_policy_id = "216adef6-5eaf-47c6-946b-b24a5dbbb65c"
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    acm_certificate_arn      = var.acm_certificate_arn
    ssl_support_method       = "sni-only"
    minimum_protocol_version = "TLSv1.2_2021"
  }

  tags = merge(var.tags, { Name = "${var.name_prefix}-api-cf" })
}
