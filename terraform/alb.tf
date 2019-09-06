resource "aws_lb" "main" {
  name               = local.name_env
  internal           = false
  load_balancer_type = "application"
  tags               = local.tags
  subnets            = aws_subnet.main.*.id
}

resource "aws_lb_listener" "main" {
  load_balancer_arn = aws_lb.main.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    target_group_arn = aws_lb_target_group.main.arn
    type             = "forward"
  }
}

resource "aws_lb_target_group" "main" {
  name        = local.name_env
  target_type = "lambda"
  vpc_id      = var.vpc_id
  tags        = local.tags
}

resource "aws_lb_target_group_attachment" "main" {
  target_group_arn = aws_lb_target_group.main.arn
  target_id        = aws_lambda_function.main.arn
}

resource "aws_lambda_permission" "main" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.main.function_name
  principal     = "elasticloadbalancing.amazonaws.com"
  source_arn    = aws_lb_target_group.main.arn
}

output "lb_url" {
  value = aws_lb.main.dns_name
}

