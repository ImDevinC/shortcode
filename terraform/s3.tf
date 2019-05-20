resource "aws_s3_bucket" "main" {
  bucket = "shortcode-${terraform.env}"
  acl    = "public-read"

  website {
    index_document = "index.html"
  }
}

output "s3_url" {
  value = "${aws_s3_bucket.main.website_endpoint}"
}
