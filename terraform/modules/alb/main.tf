resource "aws_lb" "main" {
  name               = "${var.name_prefix}-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [var.alb_sg_id]
  subnets            = var.public_subnet_ids

  enable_deletion_protection = var.environment == "production"

  # Access logs can be enabled by adding an S3 bucket here
  # access_logs { bucket = "..." enabled = true }

  tags = merge(var.tags, { Name = "${var.name_prefix}-alb" })
}

# ── Target Group for Go API (port 8080, IP-based for ECS Fargate) ──

resource "aws_lb_target_group" "api" {
  name        = "${var.name_prefix}-api-tg"
  port        = 8080
  protocol    = "HTTP"
  vpc_id      = var.vpc_id
  target_type = "ip"

  deregistration_delay = var.deregistration_delay

  health_check {
    enabled             = true
    path                = var.health_check_path
    port                = "traffic-port"
    protocol            = "HTTP"
    healthy_threshold   = 2
    unhealthy_threshold = 3
    timeout             = 5
    interval            = 30
    matcher             = "200"
  }

  tags = merge(var.tags, { Name = "${var.name_prefix}-api-tg" })

  lifecycle {
    create_before_destroy = true
  }
}

# ── HTTP listener — redirect all traffic to HTTPS ─────────────────

resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.main.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type = "redirect"
    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }
}

# ── HTTPS listener ────────────────────────────────────────────────

resource "aws_lb_listener" "https" {
  load_balancer_arn = aws_lb.main.arn
  port              = 443
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-TLS13-1-2-2021-06"
  certificate_arn   = var.acm_certificate_arn

  # Default: reject requests that don't carry the CloudFront secret header
  default_action {
    type = "fixed-response"
    fixed_response {
      content_type = "application/json"
      message_body = "{\"error\":\"Forbidden\"}"
      status_code  = "403"
    }
  }
}

# Rule: forward requests that include the correct CloudFront secret header
resource "aws_lb_listener_rule" "cloudfront_only" {
  listener_arn = aws_lb_listener.https.arn
  priority     = 10

  condition {
    http_header {
      http_header_name = var.cloudfront_secret_header_name
      values           = [var.cloudfront_secret_header_value]
    }
  }

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.api.arn
  }
}
