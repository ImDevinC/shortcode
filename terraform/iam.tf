resource "aws_iam_role" "main" {
  name               = local.tags["Name"]
  assume_role_policy = data.aws_iam_policy_document.assume_lambda.json
  tags               = local.tags
}

resource "aws_iam_role_policy_attachment" "main" {
  role       = aws_iam_role.main.name
  policy_arn = aws_iam_policy.main.arn
}

data "aws_iam_policy_document" "assume_lambda" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_iam_policy" "main" {
  name   = local.tags["Name"]
  policy = data.aws_iam_policy_document.main.json
}

data "aws_iam_policy_document" "main" {
  statement {
    sid    = "AWSLambdaBasicExecutionRole"
    effect = "Allow"

    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]

    resources = ["*"]
  }

  statement {
    sid    = "DynamodbItemCrudRole"
    effect = "Allow"

    actions = [
      "dynamodb:PutItem",
      "dynamodb:GetItem",
      "dynamodb:Query",
    ]

    resources = ["*"]
  }
}

