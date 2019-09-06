provider "aws" {
  region = "us-east-1"
}

terraform {
  backend "s3" {
    key            = "shortcode"
    encrypt        = true
    bucket         = "imdevinc-tf-storage"
    region         = "us-west-1"
    dynamodb_table = "terraform-state-lock-dynamo"
  }
}

