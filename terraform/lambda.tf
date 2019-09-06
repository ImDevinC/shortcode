resource "aws_lambda_function" "main" {
  filename         = "publish/deployment.zip"
  function_name    = local.tags["Name"]
  role             = aws_iam_role.main.arn
  handler          = "main.main"
  source_code_hash = filebase64sha256("publish/deployment.zip")
  runtime          = "python3.7"
  tags             = local.tags

  environment {
    variables = {
      DYNAMO_DB_TABLENAME = aws_dynamodb_table.main.name
    }
  }
}

