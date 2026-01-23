#!/bin/bash

echo "=========================================="
echo "Setting up DynamoDB Local for Testing"
echo "=========================================="
echo ""

# Check if DynamoDB Local container is already running
if docker ps | grep -q dynamodb-local; then
    echo "✅ DynamoDB Local is already running"
else
    echo "Starting DynamoDB Local container..."
    docker run -d \
        --name dynamodb-local \
        -p 8000:8000 \
        amazon/dynamodb-local

    if [ $? -eq 0 ]; then
        echo "✅ DynamoDB Local started successfully"
    else
        echo "❌ Failed to start DynamoDB Local"
        echo "Trying to remove old container and restart..."
        docker rm -f dynamodb-local 2>/dev/null
        docker run -d \
            --name dynamodb-local \
            -p 8000:8000 \
            amazon/dynamodb-local
    fi

    # Wait for DynamoDB to be ready
    echo "Waiting for DynamoDB Local to be ready..."
    sleep 3
fi

echo ""
echo "Creating DynamoDB tables..."
echo ""

# Create Rounds table
echo "Creating Rounds table..."
AWS_ACCESS_KEY_ID=dummy AWS_SECRET_ACCESS_KEY=dummy aws dynamodb create-table \
    --table-name AthleteUnknownRoundsDev \
    --attribute-definitions \
        AttributeName=playDate,AttributeType=S \
        AttributeName=sport,AttributeType=S \
    --key-schema \
        AttributeName=playDate,KeyType=HASH \
        AttributeName=sport,KeyType=RANGE \
    --billing-mode PAY_PER_REQUEST \
    --endpoint-url http://localhost:8000 \
    --region us-west-2 \
    2>&1 | grep -q "TableDescription" && echo "✅ Rounds table created" || echo "⚠️  Rounds table may already exist"

# Create User Stats table
echo "Creating User Stats table..."
AWS_ACCESS_KEY_ID=dummy AWS_SECRET_ACCESS_KEY=dummy aws dynamodb create-table \
    --table-name AthleteUnknownUserStatsDev \
    --attribute-definitions \
        AttributeName=userId,AttributeType=S \
    --key-schema \
        AttributeName=userId,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --endpoint-url http://localhost:8000 \
    --region us-west-2 \
    2>&1 | grep -q "TableDescription" && echo "✅ User Stats table created" || echo "⚠️  User Stats table may already exist"

echo ""
echo "=========================================="
echo "Setup Complete!"
echo "=========================================="
echo ""
echo "DynamoDB Local is running on http://localhost:8000"
echo ""
echo "To start the API server, run:"
echo ""
echo "  export DYNAMODB_ENDPOINT=http://localhost:8000"
echo "  export ROUNDS_TABLE_NAME=AthleteUnknownRoundsDev"
echo "  export USER_STATS_TABLE_NAME=AthleteUnknownUserStatsDev"
echo "  export AWS_REGION=us-west-2"
echo "  go run ."
echo ""
echo "Then in another terminal, run the tests:"
echo "  ./test-scraper-local.sh"
echo ""
echo "To stop DynamoDB Local when done:"
echo "  docker stop dynamodb-local"
echo ""
