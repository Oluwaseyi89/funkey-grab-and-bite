terraform {
  backend "s3" {
    bucket         = "funkey-terraform-state"
    key            = "staging/terraform.tfstate"
    region         = "us-east-2"
    encrypt        = true
    dynamodb_table = "funkey-terraform-locks"
  }
}
