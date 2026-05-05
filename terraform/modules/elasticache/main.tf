resource "aws_elasticache_subnet_group" "main" {
  name        = "${var.name_prefix}-elasticache-subnet-group"
  description = "ElastiCache subnet group for ${var.name_prefix}"
  subnet_ids  = var.data_subnet_ids
  tags        = var.tags
}

resource "aws_elasticache_parameter_group" "main" {
  family = "redis7"
  name   = "${var.name_prefix}-redis-params"

  parameter {
    name  = "maxmemory-policy"
    value = "allkeys-lru"
  }

  tags = var.tags
}

resource "aws_elasticache_replication_group" "main" {
  replication_group_id = "${var.name_prefix}-redis"
  description          = "Redis replication group for ${var.name_prefix} (${var.environment})"

  engine               = "redis"
  engine_version       = var.engine_version
  node_type            = var.node_type
  num_cache_clusters   = var.num_cache_clusters
  parameter_group_name = aws_elasticache_parameter_group.main.name
  subnet_group_name    = aws_elasticache_subnet_group.main.name
  security_group_ids   = [var.elasticache_sg_id]

  port = 6379

  automatic_failover_enabled = var.num_cache_clusters > 1
  multi_az_enabled           = var.num_cache_clusters > 1

  at_rest_encryption_enabled = true
  transit_encryption_enabled = true

  maintenance_window       = "sun:05:00-sun:06:00"
  snapshot_retention_limit = var.environment == "production" ? 5 : 1
  snapshot_window          = "03:00-04:00"

  apply_immediately = var.environment != "production"

  tags = merge(var.tags, { Name = "${var.name_prefix}-redis" })
}
