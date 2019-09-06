resource "aws_dynamodb_table" "main" {
  name           = "${local.tags["Name"]}"
  hash_key       = "shortcode"
  write_capacity = "${var.write_capacity}"
  read_capacity  = "${var.read_capacity}"

  attribute {
    name = "shortcode"
    type = "S"
  }

  attribute {
    name = "URI"
    type = "S"
  }

  global_secondary_index {
    name            = "URIIndex"
    hash_key        = "URI"
    projection_type = "ALL"
    write_capacity  = "${var.write_capacity}"
    read_capacity   = "${var.read_capacity}"
  }
}
