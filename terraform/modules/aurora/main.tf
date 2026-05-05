resource "random_password" "master" {
  length           = 32
  special          = true
  override_special = "!#$%&*()-_=+[]{}:?"
}

# Store the generated master password in Secrets Manager
resource "aws_secretsmanager_secret" "db_credentials" {
  name                    = "${var.name_prefix}/aurora/master-credentials"
  description             = "Aurora master credentials for ${var.name_prefix}"
  recovery_window_in_days = var.environment == "production" ? 30 : 0
  tags                    = var.tags
}

resource "aws_secretsmanager_secret_version" "db_credentials" {
  secret_id = aws_secretsmanager_secret.db_credentials.id
  secret_string = jsonencode({
    username = var.master_username
    password = random_password.master.result
    host     = aws_rds_cluster.main.endpoint
    port     = 5432
    dbname   = var.database_name
    engine   = "postgres"
  })
}

# ─────────────────────────────────────────────────────────────────
# Subnet Group & Parameter Group
# ─────────────────────────────────────────────────────────────────

resource "aws_db_subnet_group" "main" {
  name        = "${var.name_prefix}-aurora-subnet-group"
  description = "Aurora subnet group for ${var.name_prefix}"
  subnet_ids  = var.data_subnet_ids
  tags        = var.tags
}

resource "aws_rds_cluster_parameter_group" "main" {
  family      = "aurora-postgresql16"
  name        = "${var.name_prefix}-aurora-pg"
  description = "Custom parameter group for ${var.name_prefix} Aurora"

  parameter {
    name  = "log_connections"
    value = "1"
  }

  parameter {
    name  = "log_disconnections"
    value = "1"
  }

  parameter {
    name  = "log_min_duration_statement"
    value = "1000" # Log slow queries > 1 second
  }

  tags = var.tags
}

# ─────────────────────────────────────────────────────────────────
# Aurora Serverless v2 Cluster
# Note: Serverless v2 uses engine_mode = "provisioned" with
# serverlessv2_scaling_configuration block (not engine_mode = "serverless")
# ─────────────────────────────────────────────────────────────────

resource "aws_rds_cluster" "main" {
  cluster_identifier      = "${var.name_prefix}-aurora-cluster"
  engine                  = "aurora-postgresql"
  engine_mode             = "provisioned"
  engine_version          = var.engine_version
  database_name           = var.database_name
  master_username         = var.master_username
  master_password         = random_password.master.result
  db_subnet_group_name    = aws_db_subnet_group.main.name
  vpc_security_group_ids  = [var.aurora_sg_id]
  db_cluster_parameter_group_name = aws_rds_cluster_parameter_group.main.name

  serverlessv2_scaling_configuration {
    min_capacity = var.min_capacity
    max_capacity = var.max_capacity
  }

  storage_encrypted               = true
  backup_retention_period         = var.backup_retention_days
  preferred_backup_window         = "02:00-03:00"
  preferred_maintenance_window    = "sun:04:00-sun:05:00"
  deletion_protection             = var.environment == "production"
  skip_final_snapshot             = var.environment != "production"
  final_snapshot_identifier       = var.environment == "production" ? "${var.name_prefix}-aurora-final-snapshot" : null
  enabled_cloudwatch_logs_exports = ["postgresql"]

  tags = merge(var.tags, {
    Name = "${var.name_prefix}-aurora-cluster"
  })

  lifecycle {
    ignore_changes = [master_password]
  }
}

# ── Writer instance ───────────────────────────────────────────────

resource "aws_rds_cluster_instance" "writer" {
  identifier              = "${var.name_prefix}-aurora-writer"
  cluster_identifier      = aws_rds_cluster.main.id
  instance_class          = "db.serverless"
  engine                  = aws_rds_cluster.main.engine
  engine_version          = aws_rds_cluster.main.engine_version
  db_subnet_group_name    = aws_db_subnet_group.main.name
  publicly_accessible     = false
  monitoring_interval     = 60
  monitoring_role_arn     = aws_iam_role.rds_monitoring.arn
  performance_insights_enabled = true

  tags = merge(var.tags, { Name = "${var.name_prefix}-aurora-writer" })
}

# ── Reader instance (for reporting queries / Multi-AZ) ───────────

resource "aws_rds_cluster_instance" "reader" {
  count = var.enable_reader ? 1 : 0

  identifier              = "${var.name_prefix}-aurora-reader"
  cluster_identifier      = aws_rds_cluster.main.id
  instance_class          = "db.serverless"
  engine                  = aws_rds_cluster.main.engine
  engine_version          = aws_rds_cluster.main.engine_version
  db_subnet_group_name    = aws_db_subnet_group.main.name
  publicly_accessible     = false
  monitoring_interval     = 60
  monitoring_role_arn     = aws_iam_role.rds_monitoring.arn
  performance_insights_enabled = true

  tags = merge(var.tags, { Name = "${var.name_prefix}-aurora-reader" })
}

# ── Enhanced Monitoring IAM role ─────────────────────────────────

resource "aws_iam_role" "rds_monitoring" {
  name = "${var.name_prefix}-rds-monitoring-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action    = "sts:AssumeRole"
      Effect    = "Allow"
      Principal = { Service = "monitoring.rds.amazonaws.com" }
    }]
  })

  tags = var.tags
}

resource "aws_iam_role_policy_attachment" "rds_monitoring" {
  role       = aws_iam_role.rds_monitoring.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonRDSEnhancedMonitoringRole"
}
