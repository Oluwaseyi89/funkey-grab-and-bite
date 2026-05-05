# ─────────────────────────────────────────────────────────────────
# SNS Alert Topic
# ─────────────────────────────────────────────────────────────────

resource "aws_sns_topic" "alerts" {
  name              = "${var.name_prefix}-alerts"
  kms_master_key_id = "alias/aws/sns"
  tags              = var.tags
}

resource "aws_sns_topic_subscription" "email" {
  topic_arn = aws_sns_topic.alerts.arn
  protocol  = "email"
  endpoint  = var.alert_email
}

# ─────────────────────────────────────────────────────────────────
# CloudWatch Log Groups
# ─────────────────────────────────────────────────────────────────

resource "aws_cloudwatch_log_group" "api" {
  name              = "/ecs/${var.name_prefix}/api"
  retention_in_days = 30
  tags              = var.tags
}

resource "aws_cloudwatch_log_group" "workers" {
  name              = "/ecs/${var.name_prefix}/workers"
  retention_in_days = 14
  tags              = var.tags
}

resource "aws_cloudwatch_log_group" "lambda_catering" {
  name              = "/aws/lambda/${var.name_prefix}-catering-notify"
  retention_in_days = 14
  tags              = var.tags
}

# ─────────────────────────────────────────────────────────────────
# Alarms
# ─────────────────────────────────────────────────────────────────

# ECS API — high CPU
resource "aws_cloudwatch_metric_alarm" "ecs_api_cpu_high" {
  count = var.enable_ecs_api_alarms ? 1 : 0

  alarm_name          = "${var.name_prefix}-ecs-api-cpu-high"
  alarm_description   = "ECS API CPU utilisation exceeded 85%"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 3
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = 60
  statistic           = "Average"
  threshold           = 85
  treat_missing_data  = "notBreaching"

  dimensions = {
    ClusterName = var.ecs_api_cluster_name
    ServiceName = var.ecs_api_service_name
  }

  alarm_actions = [aws_sns_topic.alerts.arn]
  ok_actions    = [aws_sns_topic.alerts.arn]
}

# ECS API — high memory
resource "aws_cloudwatch_metric_alarm" "ecs_api_memory_high" {
  count = var.enable_ecs_api_alarms ? 1 : 0

  alarm_name          = "${var.name_prefix}-ecs-api-memory-high"
  alarm_description   = "ECS API memory utilisation exceeded 90%"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 3
  metric_name         = "MemoryUtilization"
  namespace           = "AWS/ECS"
  period              = 60
  statistic           = "Average"
  threshold           = 90
  treat_missing_data  = "notBreaching"

  dimensions = {
    ClusterName = var.ecs_api_cluster_name
    ServiceName = var.ecs_api_service_name
  }

  alarm_actions = [aws_sns_topic.alerts.arn]
}

# ALB — 5xx error rate
resource "aws_cloudwatch_metric_alarm" "alb_5xx" {
  count = var.enable_alb_alarms ? 1 : 0

  alarm_name          = "${var.name_prefix}-alb-5xx-errors"
  alarm_description   = "ALB 5xx error count exceeded 10 in 1 minute"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 2
  metric_name         = "HTTPCode_ELB_5XX_Count"
  namespace           = "AWS/ApplicationELB"
  period              = 60
  statistic           = "Sum"
  threshold           = 10
  treat_missing_data  = "notBreaching"

  dimensions = {
    LoadBalancer = var.alb_arn_suffix
  }

  alarm_actions = [aws_sns_topic.alerts.arn]
}

# ALB — target response time
resource "aws_cloudwatch_metric_alarm" "alb_latency" {
  count = var.enable_alb_alarms ? 1 : 0

  alarm_name          = "${var.name_prefix}-alb-latency"
  alarm_description   = "ALB p95 target response time exceeded 2 seconds"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 3
  metric_name         = "TargetResponseTime"
  namespace           = "AWS/ApplicationELB"
  period              = 60
  extended_statistic  = "p95"
  threshold           = 2
  treat_missing_data  = "notBreaching"

  dimensions = {
    LoadBalancer = var.alb_arn_suffix
  }

  alarm_actions = [aws_sns_topic.alerts.arn]
}

# Aurora — low free storage
resource "aws_cloudwatch_metric_alarm" "aurora_storage" {
  count = var.enable_aurora_alarms ? 1 : 0

  alarm_name          = "${var.name_prefix}-aurora-storage-low"
  alarm_description   = "Aurora free local storage below 2 GiB"
  comparison_operator = "LessThanThreshold"
  evaluation_periods  = 2
  metric_name         = "FreeLocalStorage"
  namespace           = "AWS/RDS"
  period              = 300
  statistic           = "Average"
  threshold           = 2147483648 # 2 GiB in bytes
  treat_missing_data  = "notBreaching"

  dimensions = {
    DBClusterIdentifier = var.aurora_cluster_identifier
  }

  alarm_actions = [aws_sns_topic.alerts.arn]
}

# Aurora — high DB connections
resource "aws_cloudwatch_metric_alarm" "aurora_connections" {
  count = var.enable_aurora_alarms ? 1 : 0

  alarm_name          = "${var.name_prefix}-aurora-connections-high"
  alarm_description   = "Aurora DB connections exceeded 50"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 3
  metric_name         = "DatabaseConnections"
  namespace           = "AWS/RDS"
  period              = 60
  statistic           = "Average"
  threshold           = 50
  treat_missing_data  = "notBreaching"

  dimensions = {
    DBClusterIdentifier = var.aurora_cluster_identifier
  }

  alarm_actions = [aws_sns_topic.alerts.arn]
}

# ElastiCache — evictions
resource "aws_cloudwatch_metric_alarm" "redis_evictions" {
  count = var.enable_redis_alarms ? 1 : 0

  alarm_name          = "${var.name_prefix}-redis-evictions"
  alarm_description   = "Redis cache evictions indicate memory pressure"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 3
  metric_name         = "Evictions"
  namespace           = "AWS/ElastiCache"
  period              = 300
  statistic           = "Sum"
  threshold           = 100
  treat_missing_data  = "notBreaching"

  dimensions = {
    ReplicationGroupId = var.elasticache_replication_group_id
  }

  alarm_actions = [aws_sns_topic.alerts.arn]
}

# ─────────────────────────────────────────────────────────────────
# CloudWatch Dashboard
# ─────────────────────────────────────────────────────────────────

resource "aws_cloudwatch_dashboard" "main" {
  dashboard_name = "${var.name_prefix}-dashboard"

  dashboard_body = jsonencode({
    widgets = [
      {
        type = "metric"
        properties = {
          title  = "ECS API — CPU & Memory"
          region = var.aws_region
          period = 60
          metrics = [
            ["AWS/ECS", "CPUUtilization", "ClusterName", var.ecs_api_cluster_name, "ServiceName", var.ecs_api_service_name, { label = "CPU %" }],
            ["AWS/ECS", "MemoryUtilization", "ClusterName", var.ecs_api_cluster_name, "ServiceName", var.ecs_api_service_name, { label = "Memory %" }]
          ]
          view = "timeSeries"
        }
      },
      {
        type = "metric"
        properties = {
          title  = "ALB — Request Count & 5xx Errors"
          region = var.aws_region
          period = 60
          metrics = [
            ["AWS/ApplicationELB", "RequestCount", "LoadBalancer", var.alb_arn_suffix, { label = "Requests" }],
            ["AWS/ApplicationELB", "HTTPCode_ELB_5XX_Count", "LoadBalancer", var.alb_arn_suffix, { label = "5xx Errors" }]
          ]
          view = "timeSeries"
        }
      },
      {
        type = "metric"
        properties = {
          title  = "Aurora — Connections & CPU"
          region = var.aws_region
          period = 60
          metrics = [
            ["AWS/RDS", "DatabaseConnections", "DBClusterIdentifier", var.aurora_cluster_identifier, { label = "Connections" }],
            ["AWS/RDS", "CPUUtilization", "DBClusterIdentifier", var.aurora_cluster_identifier, { label = "CPU %" }]
          ]
          view = "timeSeries"
        }
      }
    ]
  })
}
