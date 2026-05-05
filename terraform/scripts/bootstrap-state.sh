#!/usr/bin/env bash
# scripts/bootstrap-state.sh
# Creates the S3 bucket and DynamoDB table used for Terraform remote state.
# Run ONCE before the first 'terraform init'.
# Usage: ./scripts/bootstrap-state.sh [aws-account-id]

set -euo pipefail

ACCOUNT_ID="${1:-$(aws sts get-caller-identity --query Account --output text)}"
BUCKET="funkey-terraform-state"
TABLE="funkey-terraform-locks"
REGION="us-east-1"

echo "🔧  Bootstrapping Terraform state backend"
echo "    Bucket : ${BUCKET}"
echo "    Table  : ${TABLE}"
echo "    Region : ${REGION}"

# Create S3 bucket
aws s3api create-bucket \
  --bucket "${BUCKET}" \
  --region "${REGION}" 2>/dev/null || echo "  (bucket already exists)"

# Enable versioning
aws s3api put-bucket-versioning \
  --bucket "${BUCKET}" \
  --versioning-configuration Status=Enabled

# Enable server-side encryption
aws s3api put-bucket-encryption \
  --bucket "${BUCKET}" \
  --server-side-encryption-configuration '{
    "Rules": [{
      "ApplyServerSideEncryptionByDefault": {
        "SSEAlgorithm": "AES256"
      }
    }]
  }'

# Block all public access
aws s3api put-public-access-block \
  --bucket "${BUCKET}" \
  --public-access-block-configuration \
    "BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=true,RestrictPublicBuckets=true"

# Create DynamoDB table for state locking
aws dynamodb create-table \
  --table-name "${TABLE}" \
  --attribute-definitions AttributeName=LockID,AttributeType=S \
  --key-schema AttributeName=LockID,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST \
  --region "${REGION}" 2>/dev/null || echo "  (DynamoDB table already exists)"

echo "✅  Terraform state backend ready"
echo "    Now run: cd environments/production && terraform init"
