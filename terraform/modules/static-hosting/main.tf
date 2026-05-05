# Unique suffix so bucket name is globally unique
resource "random_id" "suffix" {
  byte_length = 4
}

locals {
  bucket_name = "${var.name_prefix}-${var.bucket_purpose}-${random_id.suffix.hex}"
}

# ─────────────────────────────────────────────────────────────────
# S3 Bucket (private — CloudFront accesses via OAC)
# ─────────────────────────────────────────────────────────────────

resource "aws_s3_bucket" "main" {
  bucket        = local.bucket_name
  force_destroy = var.force_destroy
  tags          = merge(var.tags, { Name = local.bucket_name })
}

resource "aws_s3_bucket_versioning" "main" {
  bucket = aws_s3_bucket.main.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "main" {
  bucket = aws_s3_bucket.main.id
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

# Block ALL public access — CloudFront uses OAC for private access
resource "aws_s3_bucket_public_access_block" "main" {
  bucket                  = aws_s3_bucket.main.id
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# Enforce HTTPS-only access
resource "aws_s3_bucket_policy" "https_only" {
  bucket = aws_s3_bucket.main.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid       = "DenyHTTP"
        Effect    = "Deny"
        Principal = "*"
        Action    = "s3:GetObject"
        Resource  = "${aws_s3_bucket.main.arn}/*"
        Condition = {
          Bool = { "aws:SecureTransport" = "false" }
        }
      },
      {
        Sid    = "AllowCloudFrontOAC"
        Effect = "Allow"
        Principal = {
          Service = "cloudfront.amazonaws.com"
        }
        Action   = "s3:GetObject"
        Resource = "${aws_s3_bucket.main.arn}/*"
        Condition = {
          StringEquals = {
            "AWS:SourceArn" = aws_cloudfront_origin_access_control.main.arn
          }
        }
      }
    ]
  })

  depends_on = [aws_s3_bucket_public_access_block.main]
}

# ─────────────────────────────────────────────────────────────────
# CloudFront Origin Access Control (OAC) — replaces legacy OAI
# ─────────────────────────────────────────────────────────────────

resource "aws_cloudfront_origin_access_control" "main" {
  name                              = "${local.bucket_name}-oac"
  description                       = "OAC for ${local.bucket_name}"
  origin_access_control_origin_type = "s3"
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"
}
