terraform {
  backend "s3" {
    bucket         = "funkey-terraform-state"
    key            = "production/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "funkey-terraform-locks"
  }
}
