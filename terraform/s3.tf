resource "aws_s3_bucket" "main" {
  bucket        = "shortcode-${local.tags["environment"]}"
  acl           = "public-read"
  force_destroy = true
  tags          = local.tags

  website {
    index_document = "index.html"
  }
}

output "s3_url" {
  value = aws_s3_bucket.main.website_endpoint
}

resource "aws_s3_bucket" "lambda" {
  bucket = "shortcode-lambda-${local.tags.environment}"
  acl    = "private"
  tags   = local.tags
}

resource "aws_s3_bucket_object" "lambda" {
  bucket = aws_s3_bucket.lambda.id
  key    = "${local.name_env}.zip"
  source = "publish/deployment.zip"
  etag   = filemd5("publish/deployment.zip")
}
