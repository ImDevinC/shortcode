resource "aws_lambda_function" "main" {
  filename         = "publish/deployment.zip"
  function_name    = "${local.tags["Name"]}"
  role             = "${aws_iam_role.main.arn}"
  handler          = "main"
  source_code_hash = "${base64sha256(file("publish/deployment.zip"))}"
  runtime          = "go1.x"
  tags             = "${local.tags}"

  environment {
    variables = {
      DYNAMO_DB_TABLENAME = "${aws_dynamodb_table.main.name}"
    }
  }
}

resource "aws_lambda_permission" "main" {
  action        = "lambda:InvokeFunction"
  function_name = "${local.tags["Name"]}"
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.main.execution_arn}/*/*/*"
}
