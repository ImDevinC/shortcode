variable "project" {
  default = "shortcode"
}

variable "owner" {
  default = "dcollins"
}

variable "department" {
  default = "engineering"
}

variable "sla" {
}

variable "environment" {
}

variable "write_capacity" {
}

variable "read_capacity" {
}

variable "cidr_list" {
  type = list(string)
}

variable "vpc_id" {
  default = "vpc-0a9d6870"
}


variable "log_level" {
  type    = string
  default = "debug"
}
