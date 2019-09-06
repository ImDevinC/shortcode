resource "aws_s3_bucket" "main" {
  bucket        = "shortcode-${local.tags["environment"]}"
  acl           = "public-read"
  force_destroy = true

  website {
    index_document = "index.html"
  }
}

output "s3_url" {
  value = aws_s3_bucket.main.website_endpoint
}

