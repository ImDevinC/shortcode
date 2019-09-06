locals {
  tags {
    Name        = "${local.name_env}"
    department  = "${var.department}"
    owner       = "${var.owner}"
    sla         = "${var.sla}"
    environment = "${var.environment}"
    project     = "${var.project}"
  }

  name_env = "${var.project}-${var.environment}"
}
