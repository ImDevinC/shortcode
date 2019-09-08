resource "aws_lambda_function" "main" {
  function_name    = local.tags.Name
  role             = aws_iam_role.main.arn
  handler          = "main.main"
  runtime          = "python3.7"
  s3_bucket        = aws_s3_bucket.lambda.id
  s3_key           = "${local.name_env}.zip"
  source_code_hash = aws_s3_bucket_object.lambda.etag
  tags             = local.tags

  environment {
    variables = {
      DYNAMODB_TABLE_NAME = aws_dynamodb_table.main.name
      LOG_LEVEL           = var.log_level
      HOMEPAGE            = "https://imdevinc.com"
    }
  }
}

