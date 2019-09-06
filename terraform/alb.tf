resource "aws_lb" "main" {
  name               = "${local.name_env}"
  internal           = false
  load_balancer_type = "application"
  tags               = "${local.tags}"
}
