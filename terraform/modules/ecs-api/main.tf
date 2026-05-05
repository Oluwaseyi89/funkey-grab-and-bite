data "aws_caller_identity" "current" {}

# ─────────────────────────────────────────────────────────────────
# ECS Cluster
# ─────────────────────────────────────────────────────────────────

resource "aws_ecs_cluster" "api" {
  name = "${var.name_prefix}-api-cluster"

  setting {
    name  = "containerInsights"
    value = "enabled"
  }

  tags = merge(var.tags, { Name = "${var.name_prefix}-api-cluster" })
}

resource "aws_ecs_cluster_capacity_providers" "api" {
  cluster_name = aws_ecs_cluster.api.name

  capacity_providers = ["FARGATE", "FARGATE_SPOT"]

  default_capacity_provider_strategy {
    capacity_provider = "FARGATE"
    weight            = 1
    base              = 1
  }
}

# ─────────────────────────────────────────────────────────────────
# CloudWatch Log Group
# ─────────────────────────────────────────────────────────────────

resource "aws_cloudwatch_log_group" "api" {
  name              = "/ecs/${var.name_prefix}/api"
  retention_in_days = 30
  tags              = var.tags
}

# ─────────────────────────────────────────────────────────────────
# ECS Task Definition
# ─────────────────────────────────────────────────────────────────

resource "aws_ecs_task_definition" "api" {
  family                   = "${var.name_prefix}-api"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = var.task_cpu
  memory                   = var.task_memory
  execution_role_arn       = aws_iam_role.ecs_execution.arn
  task_role_arn            = aws_iam_role.ecs_task.arn

  container_definitions = jsonencode([
    {
      name      = "funkey-api"
      image     = "${var.ecr_repository_url}:${var.image_tag}"
      essential = true

      portMappings = [
        {
          containerPort = 8080
          protocol      = "tcp"
          name          = "api"
        }
      ]

      # Pull secrets from Secrets Manager at task start — never in env vars
      secrets = [
        { name = "JWT_SECRET",          valueFrom = "${var.app_secrets_arn}:JWT_SECRET::" },
        { name = "DB_HOST",             valueFrom = "${var.db_secret_arn}:host::" },
        { name = "DB_PORT",             valueFrom = "${var.db_secret_arn}:port::" },
        { name = "DB_USER",             valueFrom = "${var.db_secret_arn}:username::" },
        { name = "DB_PASSWORD",         valueFrom = "${var.db_secret_arn}:password::" },
        { name = "DB_NAME",             valueFrom = "${var.db_secret_arn}:dbname::" },
        { name = "SES_SENDER_EMAIL",    valueFrom = "${var.app_secrets_arn}:SES_SENDER_EMAIL::" },
        { name = "TWILIO_ACCOUNT_SID",  valueFrom = "${var.app_secrets_arn}:TWILIO_ACCOUNT_SID::" },
        { name = "TWILIO_AUTH_TOKEN",   valueFrom = "${var.app_secrets_arn}:TWILIO_AUTH_TOKEN::" },
        { name = "TWILIO_PHONE_NUMBER", valueFrom = "${var.app_secrets_arn}:TWILIO_PHONE_NUMBER::" },
      ]

      environment = [
        { name = "ENVIRONMENT", value = var.environment },
        { name = "DB_SSLMODE",  value = "require" },
        { name = "AWS_REGION",  value = var.aws_region },
      ]

      healthCheck = {
        command     = ["CMD-SHELL", "curl -sf http://localhost:8080/api/v1/settings || exit 1"]
        interval    = 30
        timeout     = 5
        retries     = 3
        startPeriod = 60
      }

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.api.name
          "awslogs-region"        = var.aws_region
          "awslogs-stream-prefix" = "ecs"
        }
      }
    }
  ])

  tags = var.tags
}

# ─────────────────────────────────────────────────────────────────
# ECS Service
# ─────────────────────────────────────────────────────────────────

resource "aws_ecs_service" "api" {
  name            = "${var.name_prefix}-api-service"
  cluster         = aws_ecs_cluster.api.id
  task_definition = aws_ecs_task_definition.api.arn
  desired_count   = var.desired_count

  capacity_provider_strategy {
    capacity_provider = "FARGATE"
    weight            = 1
    base              = 1
  }

  network_configuration {
    subnets          = var.private_subnet_ids
    security_groups  = [var.ecs_api_sg_id]
    assign_public_ip = false
  }

  load_balancer {
    target_group_arn = var.target_group_arn
    container_name   = "funkey-api"
    container_port   = 8080
  }

  deployment_controller {
    type = "ECS"
  }

  deployment_circuit_breaker {
    enable   = true
    rollback = true
  }

  # Ignore task_definition changes so CI/CD can update without Terraform drift
  lifecycle {
    ignore_changes = [task_definition, desired_count]
  }

  depends_on = [aws_iam_role_policy_attachment.ecs_execution]
}

# ─────────────────────────────────────────────────────────────────
# Auto-Scaling
# ─────────────────────────────────────────────────────────────────

resource "aws_appautoscaling_target" "api" {
  max_capacity       = var.max_capacity
  min_capacity       = var.min_capacity
  resource_id        = "service/${aws_ecs_cluster.api.name}/${aws_ecs_service.api.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  service_namespace  = "ecs"
}

resource "aws_appautoscaling_policy" "api_cpu" {
  name               = "${var.name_prefix}-api-cpu-scaling"
  policy_type        = "TargetTrackingScaling"
  resource_id        = aws_appautoscaling_target.api.resource_id
  scalable_dimension = aws_appautoscaling_target.api.scalable_dimension
  service_namespace  = aws_appautoscaling_target.api.service_namespace

  target_tracking_scaling_policy_configuration {
    predefined_metric_specification {
      predefined_metric_type = "ECSServiceAverageCPUUtilization"
    }
    target_value       = 70
    scale_in_cooldown  = 300
    scale_out_cooldown = 60
  }
}

resource "aws_appautoscaling_policy" "api_memory" {
  name               = "${var.name_prefix}-api-memory-scaling"
  policy_type        = "TargetTrackingScaling"
  resource_id        = aws_appautoscaling_target.api.resource_id
  scalable_dimension = aws_appautoscaling_target.api.scalable_dimension
  service_namespace  = aws_appautoscaling_target.api.service_namespace

  target_tracking_scaling_policy_configuration {
    predefined_metric_specification {
      predefined_metric_type = "ECSServiceAverageMemoryUtilization"
    }
    target_value       = 80
    scale_in_cooldown  = 300
    scale_out_cooldown = 60
  }
}
