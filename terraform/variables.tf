variable "aws_region" {
  description = "AWS region for resources"
  type        = string
  default     = "us-west-2"
}

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
  default     = "dev"
}

variable "function_name" {
  description = "Name of the Lambda function"
  type        = string
  default     = "athlete-unknown-api-dev"
}

variable "lambda_package_path" {
  description = "Path to the Lambda deployment package"
  type        = string
  default     = "../lambda-deployment-package.zip"
}

variable "rounds_table_name" {
  description = "Name of the DynamoDB Rounds table"
  type        = string
  default     = "AthleteUnknownRoundsDev"
}

variable "user_stats_table_name" {
  description = "Name of the DynamoDB User Stats table"
  type        = string
  default     = "AthleteUnknownUserStatsDev"
}

variable "log_retention_days" {
  description = "CloudWatch log retention in days"
  type        = number
  default     = 7
}
