resource "aws_lambda_function" "lambda_function" {
  function_name = "lambda_vdlg_login"
  filename      = "./lambda.zip"
  # role          = aws_iam_role.lambda_role.arn
  role          = var.role_lambda
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  architectures = ["arm64"]
  memory_size   = 128
  timeout       = 30

  source_code_hash = filebase64sha256("./lambda.zip")

  # vpc_config {
  #   ipv6_allowed_for_dual_stack = false
  #   subnet_ids                  = data.aws_subnets.private_subnets.ids

  # }
}