variable "name_prefix" {
  description = "Resource name prefix"
  type        = string
}

variable "environment" {
  description = "Deployment environment"
  type        = string
}

variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "vpc_id" {
  description = "VPC ID"
  type        = string
}

variable "private_subnet_ids" {
  description = "Private subnet IDs for the Lambda function"
  type        = list(string)
}

variable "lambda_sg_id" {
  description = "Security group ID for Lambda functions"
  type        = string
}

variable "ecr_repository_url" {
  description = "ECR repository URL for the Lambda container image"
  type        = string
}

variable "image_tag" {
  description = "Container image tag"
  type        = string
  default     = "latest"
}

variable "sqs_catering_queue_arn" {
  description = "SQS catering queue ARN (event source)"
  type        = string
}

variable "sqs_catering_queue_url" {
  description = "SQS catering queue URL"
  type        = string
}

variable "app_secrets_arn" {
  description = "Secrets Manager ARN containing SES and Twilio credentials"
  type        = string
}

variable "ses_sender_email" {
  description = "Verified SES sender email"
  type        = string
  default     = "noreply@funkeygrabandbite.com"
}

variable "reserved_concurrency" {
  description = "Lambda reserved concurrency (-1 = unreserved)"
  type        = number
  default     = 5
}

variable "tags" {
  description = "Tags applied to all resources"
  type        = map(string)
  default     = {}
}
