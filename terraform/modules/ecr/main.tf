data "aws_caller_identity" "current" {}

locals {
  lifecycle_policy = jsonencode({
    rules = [
      {
        rulePriority = 1
        description  = "Keep last ${var.image_count_to_keep} tagged images"
        selection = {
          tagStatus   = "tagged"
          tagPrefixList = ["v", "sha-", "latest"]
          countType   = "imageCountMoreThan"
          countNumber = var.image_count_to_keep
        }
        action = { type = "expire" }
      },
      {
        rulePriority = 2
        description  = "Expire untagged images older than 7 days"
        selection = {
          tagStatus   = "untagged"
          countType   = "sinceImagePushed"
          countUnit   = "days"
          countNumber = 7
        }
        action = { type = "expire" }
      }
    ]
  })
}

# ── API repository ────────────────────────────────────────────────

resource "aws_ecr_repository" "api" {
  name                 = "${var.name_prefix}/api"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  encryption_configuration {
    encryption_type = "AES256"
  }

  tags = merge(var.tags, { Name = "${var.name_prefix}-api-ecr" })
}

resource "aws_ecr_lifecycle_policy" "api" {
  repository = aws_ecr_repository.api.name
  policy     = local.lifecycle_policy
}

# ── Lambda catering-notify repository ────────────────────────────

resource "aws_ecr_repository" "lambda_catering" {
  name                 = "${var.name_prefix}/lambda-catering"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  encryption_configuration {
    encryption_type = "AES256"
  }

  tags = merge(var.tags, { Name = "${var.name_prefix}-lambda-catering-ecr" })
}

resource "aws_ecr_lifecycle_policy" "lambda_catering" {
  repository = aws_ecr_repository.lambda_catering.name
  policy     = local.lifecycle_policy
}
