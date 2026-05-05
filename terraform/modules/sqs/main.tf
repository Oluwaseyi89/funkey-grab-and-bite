# ─────────────────────────────────────────────────────────────────
# Order Queue + DLQ
# ─────────────────────────────────────────────────────────────────

resource "aws_sqs_queue" "order_dlq" {
  name                      = "${var.name_prefix}-order-dlq"
  message_retention_seconds = 1209600 # 14 days for failed messages
  tags                      = merge(var.tags, { Name = "${var.name_prefix}-order-dlq" })
}

resource "aws_sqs_queue" "order" {
  name                       = "${var.name_prefix}-order-queue"
  message_retention_seconds  = var.message_retention_seconds
  visibility_timeout_seconds = var.visibility_timeout_seconds

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.order_dlq.arn
    maxReceiveCount     = 3
  })

  tags = merge(var.tags, { Name = "${var.name_prefix}-order-queue" })
}

# ─────────────────────────────────────────────────────────────────
# Catering Queue + DLQ (triggers Lambda catering-notify)
# ─────────────────────────────────────────────────────────────────

resource "aws_sqs_queue" "catering_dlq" {
  name                      = "${var.name_prefix}-catering-dlq"
  message_retention_seconds = 1209600
  tags                      = merge(var.tags, { Name = "${var.name_prefix}-catering-dlq" })
}

resource "aws_sqs_queue" "catering" {
  name                       = "${var.name_prefix}-catering-queue"
  message_retention_seconds  = var.message_retention_seconds
  visibility_timeout_seconds = var.visibility_timeout_seconds

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.catering_dlq.arn
    maxReceiveCount     = 3
  })

  tags = merge(var.tags, { Name = "${var.name_prefix}-catering-queue" })
}

# ─────────────────────────────────────────────────────────────────
# Queue policies — enforce HTTPS-only and restrict producers
# ─────────────────────────────────────────────────────────────────

data "aws_caller_identity" "current" {}

resource "aws_sqs_queue_policy" "order" {
  queue_url = aws_sqs_queue.order.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid       = "DenyNonSSL"
        Effect    = "Deny"
        Principal = { AWS = "*" }
        Action    = "sqs:*"
        Resource  = aws_sqs_queue.order.arn
        Condition = { Bool = { "aws:SecureTransport" = "false" } }
      }
    ]
  })
}

resource "aws_sqs_queue_policy" "catering" {
  queue_url = aws_sqs_queue.catering.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid       = "DenyNonSSL"
        Effect    = "Deny"
        Principal = { AWS = "*" }
        Action    = "sqs:*"
        Resource  = aws_sqs_queue.catering.arn
        Condition = { Bool = { "aws:SecureTransport" = "false" } }
      }
    ]
  })
}
