resource "aws_api_gateway_rest_api" "main" {
  name        = "${local.tags["Name"]}"
  description = "Generates and retrieves shortcodes"
}

resource "aws_api_gateway_resource" "main" {
  rest_api_id = "${aws_api_gateway_rest_api.main.id}"
  parent_id   = "${aws_api_gateway_rest_api.main.root_resource_id}"
  path_part   = "{PROXY+}"
}

resource "aws_api_gateway_method" "get" {
  rest_api_id   = "${aws_api_gateway_rest_api.main.id}"
  resource_id   = "${aws_api_gateway_resource.main.id}"
  http_method   = "GET"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "get" {
  rest_api_id             = "${aws_api_gateway_rest_api.main.id}"
  resource_id             = "${aws_api_gateway_method.get.resource_id}"
  http_method             = "${aws_api_gateway_method.get.http_method}"
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = "${aws_lambda_function.main.invoke_arn}"
}

resource "aws_api_gateway_method" "post" {
  rest_api_id   = "${aws_api_gateway_rest_api.main.id}"
  resource_id   = "${aws_api_gateway_resource.main.id}"
  http_method   = "POST"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "post" {
  rest_api_id             = "${aws_api_gateway_rest_api.main.id}"
  resource_id             = "${aws_api_gateway_method.post.resource_id}"
  http_method             = "${aws_api_gateway_method.post.http_method}"
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = "${aws_lambda_function.main.invoke_arn}"
}

resource "aws_api_gateway_deployment" "main" {
  rest_api_id = "${aws_api_gateway_rest_api.main.id}"
  stage_name  = "main"
}

output "url" {
  value = "${aws_api_gateway_deployment.main.invoke_url}"
}
