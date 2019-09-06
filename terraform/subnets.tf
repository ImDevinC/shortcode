resource "aws_subnet" "main" {
  count      = length(var.cidr_list)
  vpc_id     = var.vpc_id
  cidr_block = element(var.cidr_list, count.index)
}

