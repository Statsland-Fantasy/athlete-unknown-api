#!/bin/bash

# Build script for AWS Lambda deployment
# This script builds the Go application for Lambda and creates a deployment package

set -e

echo "Building Lambda function for athlete-unknown-api..."

# Set build variables
BINARY_NAME="bootstrap"
PACKAGE_NAME="lambda-deployment-package.zip"

# Clean previous builds
echo "Cleaning previous builds..."
rm -f ${BINARY_NAME}
rm -f ${PACKAGE_NAME}

# Build for Linux AMD64 (Lambda execution environment)
echo "Building Go binary for Linux AMD64..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda -ldflags="-s -w" -o ${BINARY_NAME} .

# Create deployment package
echo "Creating deployment package..."
zip ${PACKAGE_NAME} ${BINARY_NAME}

# Display package info
echo ""
echo "Build complete!"
echo "Binary: ${BINARY_NAME}"
echo "Package: ${PACKAGE_NAME}"
echo "Package size: $(du -h ${PACKAGE_NAME} | cut -f1)"
echo ""
echo "To deploy to AWS Lambda, upload ${PACKAGE_NAME} to your Lambda function."
