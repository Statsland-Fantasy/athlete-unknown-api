terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Environment = var.environment
      Project     = "AthleteUnknown"
      ManagedBy   = "Terraform"
    }
  }
}

# IAM Role for Lambda
resource "aws_iam_role" "lambda_role" {
  name = "${var.function_name}-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

# Attach basic Lambda execution policy
resource "aws_iam_role_policy_attachment" "lambda_basic" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# IAM Policy for DynamoDB access
resource "aws_iam_role_policy" "lambda_dynamodb" {
  name = "${var.function_name}-dynamodb-policy"
  role = aws_iam_role.lambda_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem",
          "dynamodb:DeleteItem",
          "dynamodb:Scan",
          "dynamodb:Query",
          "dynamodb:BatchGetItem",
          "dynamodb:BatchWriteItem"
        ]
        Resource = [
          "arn:aws:dynamodb:${var.aws_region}:*:table/${var.rounds_table_name}",
          "arn:aws:dynamodb:${var.aws_region}:*:table/${var.user_stats_table_name}"
        ]
      }
    ]
  })
}

# Lambda Function
resource "aws_lambda_function" "api" {
  filename         = var.lambda_package_path
  function_name    = var.function_name
  role            = aws_iam_role.lambda_role.arn
  handler         = "bootstrap"
  source_code_hash = filebase64sha256(var.lambda_package_path)
  runtime         = "provided.al2023"
  timeout         = 30
  memory_size     = 512

  environment {
    variables = {
      ROUNDS_TABLE_NAME    = var.rounds_table_name
      USER_STATS_TABLE_NAME = var.user_stats_table_name
      AWS_REGION           = var.aws_region
    }
  }
}

# CloudWatch Log Group
resource "aws_cloudwatch_log_group" "lambda_logs" {
  name              = "/aws/lambda/${var.function_name}"
  retention_in_days = var.log_retention_days
}

# Lambda Function URL (for easier testing, can be replaced with API Gateway)
resource "aws_lambda_function_url" "api_url" {
  function_name      = aws_lambda_function.api.function_name
  authorization_type = "NONE"

  cors {
    allow_origins     = ["*"]
    allow_methods     = ["GET", "POST", "DELETE", "OPTIONS"]
    allow_headers     = ["Content-Type", "Authorization", "X-API-Key"]
    max_age          = 300
  }
}

# Outputs
output "lambda_function_name" {
  description = "Name of the Lambda function"
  value       = aws_lambda_function.api.function_name
}

output "lambda_function_arn" {
  description = "ARN of the Lambda function"
  value       = aws_lambda_function.api.arn
}

output "lambda_role_arn" {
  description = "ARN of the Lambda IAM role"
  value       = aws_iam_role.lambda_role.arn
}

output "function_url" {
  description = "URL to invoke the Lambda function"
  value       = aws_lambda_function_url.api_url.function_url
}

output "cloudwatch_log_group" {
  description = "CloudWatch Log Group name"
  value       = aws_cloudwatch_log_group.lambda_logs.name
}
