resource "aws_ecs_cluster" "workers" {
  name = "${var.name_prefix}-workers-cluster"

  setting {
    name  = "containerInsights"
    value = "enabled"
  }

  tags = merge(var.tags, { Name = "${var.name_prefix}-workers-cluster" })
}

resource "aws_ecs_cluster_capacity_providers" "workers" {
  cluster_name       = aws_ecs_cluster.workers.name
  capacity_providers = ["FARGATE", "FARGATE_SPOT"]

  # Prefer FARGATE_SPOT for cost savings on background workers
  default_capacity_provider_strategy {
    capacity_provider = "FARGATE_SPOT"
    weight            = 4
    base              = 0
  }

  default_capacity_provider_strategy {
    capacity_provider = "FARGATE"
    weight            = 1
    base              = 1
  }
}

resource "aws_cloudwatch_log_group" "workers" {
  name              = "/ecs/${var.name_prefix}/workers"
  retention_in_days = 14
  tags              = var.tags
}

resource "aws_ecs_task_definition" "workers" {
  family                   = "${var.name_prefix}-workers"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = var.task_cpu
  memory                   = var.task_memory
  execution_role_arn       = aws_iam_role.ecs_execution.arn
  task_role_arn            = aws_iam_role.ecs_task.arn

  container_definitions = jsonencode([
    {
      name      = "funkey-worker"
      image     = "${var.ecr_repository_url}:${var.image_tag}"
      essential = true

      command = ["./funkey-bite-api", "-mode=worker"]

      secrets = [
        { name = "DB_HOST", valueFrom = "${var.db_secret_arn}:host::" },
        { name = "DB_PORT", valueFrom = "${var.db_secret_arn}:port::" },
        { name = "DB_USER", valueFrom = "${var.db_secret_arn}:username::" },
        { name = "DB_PASSWORD", valueFrom = "${var.db_secret_arn}:password::" },
        { name = "DB_NAME", valueFrom = "${var.db_secret_arn}:dbname::" },
        { name = "JWT_SECRET", valueFrom = "${var.app_secrets_arn}:JWT_SECRET::" },
      ]

      environment = [
        { name = "ENVIRONMENT", value = var.environment },
        { name = "DB_SSLMODE", value = "require" },
        { name = "AWS_REGION", value = var.aws_region },
        { name = "SQS_ORDER_QUEUE_URL", value = var.sqs_order_queue_url },
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.workers.name
          "awslogs-region"        = var.aws_region
          "awslogs-stream-prefix" = "ecs"
        }
      }
    }
  ])

  tags = var.tags
}

resource "aws_ecs_service" "workers" {
  name            = "${var.name_prefix}-workers-service"
  cluster         = aws_ecs_cluster.workers.id
  task_definition = aws_ecs_task_definition.workers.arn
  desired_count   = var.desired_count

  capacity_provider_strategy {
    capacity_provider = "FARGATE_SPOT"
    weight            = 4
    base              = 0
  }

  capacity_provider_strategy {
    capacity_provider = "FARGATE"
    weight            = 1
    base              = 1
  }

  network_configuration {
    subnets          = var.private_subnet_ids
    security_groups  = [var.ecs_workers_sg_id]
    assign_public_ip = false
  }

  deployment_circuit_breaker {
    enable   = true
    rollback = true
  }

  lifecycle {
    ignore_changes = [task_definition, desired_count]
  }

  depends_on = [aws_iam_role_policy_attachment.ecs_execution]
}
