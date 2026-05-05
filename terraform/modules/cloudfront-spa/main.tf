# CloudFront distribution for Single-Page Applications (Nuxt web + React admin)
# Features:
#   - S3 origin accessed via OAC (no public S3 access)
#   - SPA routing: 403/404 → /index.html (HTTP 200)
#   - HTTPS redirect enforced
#   - Optimised caching via AWS managed policy
#   - WAF ACL attached when provided

locals {
  origin_id = "${var.name_prefix}-${var.app_name}-s3-origin"
}

resource "aws_cloudfront_distribution" "main" {
  enabled             = true
  is_ipv6_enabled     = true
  default_root_object = "index.html"
  price_class         = var.price_class
  aliases             = var.domain_aliases
  web_acl_id          = var.waf_web_acl_arn != "" ? var.waf_web_acl_arn : null
  comment             = "Funkey ${var.app_name} — ${var.environment}"

  origin {
    domain_name              = var.bucket_regional_domain_name
    origin_id                = local.origin_id
    origin_access_control_id = var.oac_id
  }

  default_cache_behavior {
    target_origin_id       = local.origin_id
    viewer_protocol_policy = "redirect-to-https"
    allowed_methods        = ["GET", "HEAD", "OPTIONS"]
    cached_methods         = ["GET", "HEAD"]
    compress               = true

    # AWS managed CachingOptimized policy
    cache_policy_id = "658327ea-f89d-4fab-a63d-7e88639e58f6"
  }

  # SPA routing: redirect all 403/404 to index.html so the client router handles it
  custom_error_response {
    error_code            = 403
    response_code         = 200
    response_page_path    = "/index.html"
    error_caching_min_ttl = 0
  }

  custom_error_response {
    error_code            = 404
    response_code         = 200
    response_page_path    = "/index.html"
    error_caching_min_ttl = 0
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

  tags = merge(var.tags, { Name = "${var.name_prefix}-${var.app_name}-cf" })
}
