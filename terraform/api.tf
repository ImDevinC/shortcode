# resource "aws_api_gateway_rest_api" "main" {
#   name        = "${local.tags["Name"]}"
#   description = "Generates and retrieves shortcodes"
# }
# resource "aws_api_gateway_resource" "proxy" {
#   rest_api_id = "${aws_api_gateway_rest_api.main.id}"
#   parent_id   = "${aws_api_gateway_rest_api.main.root_resource_id}"
#   path_part   = "{PROXY+}"
# }
# #########################
# ## Root Gateway Method ##
# #########################
# resource "aws_api_gateway_method" "root_get" {
#   rest_api_id   = "${aws_api_gateway_rest_api.main.id}"
#   resource_id   = "${aws_api_gateway_rest_api.main.root_resource_id}"
#   http_method   = "GET"
#   authorization = "NONE"
# }
# resource "aws_api_gateway_integration" "root_get" {
#   rest_api_id             = "${aws_api_gateway_rest_api.main.id}"
#   resource_id             = "${aws_api_gateway_method.root_get.resource_id}"
#   http_method             = "${aws_api_gateway_method.root_get.http_method}"
#   integration_http_method = "POST"
#   type                    = "AWS_PROXY"
#   uri                     = "${aws_lambda_function.main.invoke_arn}"
# }
# resource "aws_api_gateway_method" "root_post" {
#   rest_api_id   = "${aws_api_gateway_rest_api.main.id}"
#   resource_id   = "${aws_api_gateway_rest_api.main.root_resource_id}"
#   http_method   = "POST"
#   authorization = "NONE"
# }
# resource "aws_api_gateway_integration" "root_post" {
#   rest_api_id             = "${aws_api_gateway_rest_api.main.id}"
#   resource_id             = "${aws_api_gateway_method.root_post.resource_id}"
#   http_method             = "${aws_api_gateway_method.root_post.http_method}"
#   integration_http_method = "POST"
#   type                    = "AWS_PROXY"
#   uri                     = "${aws_lambda_function.main.invoke_arn}"
# }
# #############################
# ## {PROXY+} Gateway Method ##
# #############################
# resource "aws_api_gateway_method" "proxy_get" {
#   rest_api_id   = "${aws_api_gateway_rest_api.main.id}"
#   resource_id   = "${aws_api_gateway_resource.proxy.id}"
#   http_method   = "GET"
#   authorization = "NONE"
# }
# resource "aws_api_gateway_integration" "proxy_get" {
#   rest_api_id             = "${aws_api_gateway_rest_api.main.id}"
#   resource_id             = "${aws_api_gateway_method.proxy_get.resource_id}"
#   http_method             = "${aws_api_gateway_method.proxy_get.http_method}"
#   integration_http_method = "POST"
#   type                    = "AWS_PROXY"
#   uri                     = "${aws_lambda_function.main.invoke_arn}"
# }
# resource "aws_api_gateway_method" "proxy_post" {
#   rest_api_id   = "${aws_api_gateway_rest_api.main.id}"
#   resource_id   = "${aws_api_gateway_resource.proxy.id}"
#   http_method   = "POST"
#   authorization = "NONE"
# }
# resource "aws_api_gateway_integration" "proxy_post" {
#   rest_api_id             = "${aws_api_gateway_rest_api.main.id}"
#   resource_id             = "${aws_api_gateway_method.proxy_post.resource_id}"
#   http_method             = "${aws_api_gateway_method.proxy_post.http_method}"
#   integration_http_method = "POST"
#   type                    = "AWS_PROXY"
#   uri                     = "${aws_lambda_function.main.invoke_arn}"
# }
# #############################
# ##         Stages          ##
# #############################
# resource "aws_api_gateway_deployment" "main" {
#   rest_api_id = "${aws_api_gateway_rest_api.main.id}"
#   stage_name  = "main"
# }
# output "url" {
#   value = "${aws_api_gateway_deployment.main.invoke_url}"
# }

