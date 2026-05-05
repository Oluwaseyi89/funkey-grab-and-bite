# ─────────────────────────────────────────────────────────────────
# Custom Event Bus
# ─────────────────────────────────────────────────────────────────

resource "aws_cloudwatch_event_bus" "main" {
  name = "${var.name_prefix}-event-bus"
  tags = var.tags
}

# ─────────────────────────────────────────────────────────────────
# Scheduler IAM Role
# ─────────────────────────────────────────────────────────────────

data "aws_iam_policy_document" "scheduler_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["scheduler.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "scheduler" {
  name               = "${var.name_prefix}-eventbridge-scheduler-role"
  assume_role_policy = data.aws_iam_policy_document.scheduler_assume_role.json
  tags               = var.tags
}

resource "aws_iam_role_policy" "scheduler_ecs" {
  name = "${var.name_prefix}-scheduler-ecs-policy"
  role = aws_iam_role.scheduler.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "ecs:RunTask",
          "ecs:DescribeTasks"
        ]
        Resource = "*"
      },
      {
        Effect = "Allow"
        Action = ["iam:PassRole"]
        Resource = [
          var.ecs_task_role_arn != "" ? var.ecs_task_role_arn : "*"
        ]
      }
    ]
  })
}

# ─────────────────────────────────────────────────────────────────
# Scheduled Rules
# ─────────────────────────────────────────────────────────────────

# Daily promotion expiry check — runs nightly at 00:05 UTC
resource "aws_scheduler_schedule" "promotion_expiry" {
  count = var.ecs_api_cluster_arn != "" ? 1 : 0

  name       = "${var.name_prefix}-promotion-expiry"
  group_name = "default"
  description = "Expire promotions that have passed valid_until"

  schedule_expression          = "cron(5 0 * * ? *)"
  schedule_expression_timezone = "UTC"

  flexible_time_window {
    mode = "OFF"
  }

  target {
    arn      = var.ecs_api_cluster_arn
    role_arn = aws_iam_role.scheduler.arn

    ecs_parameters {
      task_definition_arn = var.ecs_api_task_definition_arn
      task_count          = 1
      launch_type         = "FARGATE"

      network_configuration {
        assign_public_ip = false
        security_groups  = [var.ecs_api_sg_id]
        subnets          = var.private_subnet_ids
      }

      overrides {
        container_override {
          name    = "funkey-api"
          command = ["./funkey-bite-api", "-task=expire-promotions"]
        }
      }
    }
  }
}

# Daily sales report snapshot — runs at 01:00 UTC
resource "aws_scheduler_schedule" "daily_report" {
  count = var.ecs_api_cluster_arn != "" ? 1 : 0

  name        = "${var.name_prefix}-daily-report"
  group_name  = "default"
  description = "Generate and store daily sales snapshot"

  schedule_expression          = "cron(0 1 * * ? *)"
  schedule_expression_timezone = "UTC"

  flexible_time_window {
    mode                      = "FLEXIBLE"
    maximum_window_in_minutes = 15
  }

  target {
    arn      = var.ecs_api_cluster_arn
    role_arn = aws_iam_role.scheduler.arn

    ecs_parameters {
      task_definition_arn = var.ecs_api_task_definition_arn
      task_count          = 1
      launch_type         = "FARGATE"

      network_configuration {
        assign_public_ip = false
        security_groups  = [var.ecs_api_sg_id]
        subnets          = var.private_subnet_ids
      }

      overrides {
        container_override {
          name    = "funkey-api"
          command = ["./funkey-bite-api", "-task=daily-report"]
        }
      }
    }
  }
}
