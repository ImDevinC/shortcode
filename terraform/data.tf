locals {
  tags {
    Name        = "${var.project}-${var.environment}"
    department  = "${var.department}"
    owner       = "${var.owner}"
    sla         = "${var.sla}"
    environment = "${var.environment}"
    project     = "${var.project}"
  }
}
