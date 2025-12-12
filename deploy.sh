#!/bin/bash

# Automated deployment script for AWS Lambda
# This script builds the Lambda function and deploys it using Terraform

set -e

echo "=========================================="
echo "Athlete Unknown API - Lambda Deployment"
echo "=========================================="
echo ""

# Step 1: Build Lambda function
echo "Step 1: Building Lambda function..."
./build-lambda.sh

if [ ! -f "lambda-deployment-package.zip" ]; then
    echo "Error: Lambda package not found. Build failed."
    exit 1
fi

echo ""
echo "Step 2: Deploying with Terraform..."
cd terraform

# Check if terraform is initialized
if [ ! -d ".terraform" ]; then
    echo "Initializing Terraform..."
    terraform init
fi

# Check if tfvars file exists
if [ ! -f "terraform.tfvars" ]; then
    echo ""
    echo "Warning: terraform.tfvars not found!"
    echo "Please create terraform.tfvars from terraform.tfvars.example"
    echo "and configure your variables before deploying."
    echo ""
    read -p "Continue with default values? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Plan
echo ""
echo "Generating Terraform plan..."
terraform plan -out=tfplan

# Confirm
echo ""
read -p "Apply this plan? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Deployment cancelled."
    rm tfplan
    exit 0
fi

# Apply
echo ""
echo "Applying Terraform configuration..."
terraform apply tfplan
rm tfplan

echo ""
echo "=========================================="
echo "Deployment Complete!"
echo "=========================================="
echo ""
echo "To test your API, use the function_url output from Terraform:"
echo ""
terraform output function_url
echo ""
echo "Example test:"
echo '  curl "$(terraform output -raw function_url)/health"'
echo ""
