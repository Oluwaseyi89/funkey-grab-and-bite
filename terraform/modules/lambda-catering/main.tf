resource "aws_cloudwatch_log_group" "lambda" {
  name              = "/aws/lambda/${var.name_prefix}-catering-notify"
  retention_in_days = 14
  tags              = var.tags
}

# ─────────────────────────────────────────────────────────────────
# Lambda function (container image)
# CI/CD pipeline builds and pushes the image; Terraform manages config
# ─────────────────────────────────────────────────────────────────

resource "aws_lambda_function" "catering_notify" {
  function_name = "${var.name_prefix}-catering-notify"
  description   = "Processes catering request events and sends email/SMS notifications"
  role          = aws_iam_role.lambda.arn
  package_type  = "Image"
  image_uri     = "${var.ecr_repository_url}:${var.image_tag}"
  timeout       = 30
  memory_size   = 256

  reserved_concurrent_executions = var.reserved_concurrency

  # Deploy inside VPC so it can reach the same Secrets Manager endpoint
  vpc_config {
    subnet_ids         = var.private_subnet_ids
    security_group_ids = [var.lambda_sg_id]
  }

  environment {
    variables = {
      ENVIRONMENT      = var.environment
      AWS_REGION_NAME  = var.aws_region
      SES_SENDER_EMAIL = var.ses_sender_email
      SECRETS_ARN      = var.app_secrets_arn
      SQS_QUEUE_URL    = var.sqs_catering_queue_url
    }
  }

  depends_on = [
    aws_iam_role_policy_attachment.lambda_basic,
    aws_cloudwatch_log_group.lambda,
  ]

  tags = merge(var.tags, { Name = "${var.name_prefix}-catering-notify" })
}

# ─────────────────────────────────────────────────────────────────
# SQS Event Source Mapping
# ─────────────────────────────────────────────────────────────────

resource "aws_lambda_event_source_mapping" "catering_sqs" {
  event_source_arn                   = var.sqs_catering_queue_arn
  function_name                      = aws_lambda_function.catering_notify.arn
  batch_size                         = 5
  maximum_batching_window_in_seconds = 10
  enabled                            = true

  function_response_types = ["ReportBatchItemFailures"]
}
