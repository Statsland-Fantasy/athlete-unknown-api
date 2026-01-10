package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DB wraps the DynamoDB client
type DB struct {
	client             *dynamodb.Client
	roundsTableName    string
	userStatsTableName string
}

// NewDB creates a new DynamoDB client
func NewDB(cfg *Config) (*DB, error) {
	var awsCfg aws.Config
	var err error

	// Load AWS config
	if cfg.DynamoDBEndpoint != "" {
		// For local DynamoDB or custom endpoint - use dummy credentials
		awsCfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(cfg.AWSRegion),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "")),
		)
		if err != nil {
			return nil, fmt.Errorf("unable to load AWS SDK config: %w", err)
		}

		// Create client with custom endpoint
		client := dynamodb.NewFromConfig(awsCfg, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(cfg.DynamoDBEndpoint)
		})

		return &DB{
			client:             client,
			roundsTableName:    cfg.RoundsTableName,
			userStatsTableName: cfg.UserStatsTableName,
		}, nil
	}

	// For standard AWS DynamoDB - use real credentials
	awsCfg, err = config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.AWSRegion),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS SDK config: %w", err)
	}

	client := dynamodb.NewFromConfig(awsCfg)

	return &DB{
		client:             client,
		roundsTableName:    cfg.RoundsTableName,
		userStatsTableName: cfg.UserStatsTableName,
	}, nil
}

// GetRound retrieves a round by playDate and sport
func (db *DB) GetRound(ctx context.Context, sport, playDate string) (*Round, error) {
	result, err := db.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(db.roundsTableName),
		Key: map[string]types.AttributeValue{
			"playDate": &types.AttributeValueMemberS{Value: playDate},
			"sport":    &types.AttributeValueMemberS{Value: sport},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get round: %w", err)
	}

	if result.Item == nil {
		return nil, nil // Not found
	}

	var round Round
	err = attributevalue.UnmarshalMap(result.Item, &round)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal round: %w", err)
	}

	return &round, nil
}

// CreateRound creates a new round
func (db *DB) CreateRound(ctx context.Context, round *Round) error {
	// Set timestamps
	now := time.Now()
	round.Created = now
	round.LastUpdated = now

	// Initialize stats if not provided
	if round.Stats.PlayDate == "" {
		round.Stats.PlayDate = round.PlayDate
		round.Stats.Name = round.Player.Name
		round.Stats.Sport = round.Sport
	}

	// Marshal the round to DynamoDB format
	item, err := attributevalue.MarshalMap(round)
	if err != nil {
		return fmt.Errorf("failed to marshal round: %w", err)
	}

	// Check if item already exists using ConditionExpression
	_, err = db.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(db.roundsTableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(playDate) AND attribute_not_exists(sport)"),
	})
	if err != nil {
		if _, ok := err.(*types.ConditionalCheckFailedException); ok {
			return fmt.Errorf("round already exists")
		}
		return fmt.Errorf("failed to create round: %w", err)
	}

	return nil
}

// UpdateRound updates an existing round
func (db *DB) UpdateRound(ctx context.Context, round *Round) error {
	round.LastUpdated = time.Now()

	// Marshal the round to DynamoDB format
	item, err := attributevalue.MarshalMap(round)
	if err != nil {
		return fmt.Errorf("failed to marshal round: %w", err)
	}

	_, err = db.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(db.roundsTableName),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to update round: %w", err)
	}

	return nil
}

// DeleteRound deletes a round by playDate and sport
func (db *DB) DeleteRound(ctx context.Context, sport, playDate string) error {
	_, err := db.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(db.roundsTableName),
		Key: map[string]types.AttributeValue{
			"playDate": &types.AttributeValueMemberS{Value: playDate},
			"sport":    &types.AttributeValueMemberS{Value: sport},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete round: %w", err)
	}

	return nil
}

// GetRoundsBySport retrieves minimal round information for a specific sport, optionally filtered by date range
// Returns only roundId, sport, and playDate fields using DynamoDB ProjectionExpression for efficiency
// Uses the SportPlayDateIndex GSI for efficient querying and automatic sorting by playDate
func (db *DB) GetRoundsBySport(ctx context.Context, sport, startDate, endDate string) ([]*RoundSummary, error) {
	// Build key condition expression for sport (partition key of GSI)
	keyConditionExpression := "sport = :sport"
	expressionAttributeValues := map[string]types.AttributeValue{
		":sport": &types.AttributeValueMemberS{Value: sport},
	}

	// Add date range filtering to key condition if provided
	if startDate != "" && endDate != "" {
		keyConditionExpression += " AND playDate BETWEEN :startDate AND :endDate"
		expressionAttributeValues[":startDate"] = &types.AttributeValueMemberS{Value: startDate}
		expressionAttributeValues[":endDate"] = &types.AttributeValueMemberS{Value: endDate}
	} else if startDate != "" {
		keyConditionExpression += " AND playDate >= :startDate"
		expressionAttributeValues[":startDate"] = &types.AttributeValueMemberS{Value: startDate}
	} else if endDate != "" {
		keyConditionExpression += " AND playDate <= :endDate"
		expressionAttributeValues[":endDate"] = &types.AttributeValueMemberS{Value: endDate}
	}

	result, err := db.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(db.roundsTableName),
		IndexName:                 aws.String("SportPlayDateIndex"),
		KeyConditionExpression:    aws.String(keyConditionExpression),
		ExpressionAttributeValues: expressionAttributeValues,
		ProjectionExpression:      aws.String("roundId, sport, playDate"),
		ScanIndexForward:          aws.Bool(false), // Sort descending (latest to earliest)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query rounds: %w", err)
	}

	var rounds []*RoundSummary
	for _, item := range result.Items {
		var round RoundSummary
		err = attributevalue.UnmarshalMap(item, &round)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal round: %w", err)
		}
		rounds = append(rounds, &round)
	}

	return rounds, nil
}

// GetUserStats retrieves user statistics by userId
func (db *DB) GetUserStats(ctx context.Context, userId string) (*UserStats, error) {
	result, err := db.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(db.userStatsTableName),
		Key: map[string]types.AttributeValue{
			ConstantUserId: &types.AttributeValueMemberS{Value: userId},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}

	if result.Item == nil {
		return nil, nil // Not found
	}

	var stats UserStats
	err = attributevalue.UnmarshalMap(result.Item, &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user stats: %w", err)
	}

	return &stats, nil
}

// CreateUserStats creates new user statistics
func (db *DB) CreateUserStats(ctx context.Context, stats *UserStats) error {
	// Set timestamp
	stats.UserCreated = time.Now()

	// Marshal the stats to DynamoDB format
	item, err := attributevalue.MarshalMap(stats)
	if err != nil {
		return fmt.Errorf("failed to marshal user stats: %w", err)
	}

	// Check if item already exists using ConditionExpression
	_, err = db.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(db.userStatsTableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(userId)"),
	})
	if err != nil {
		if _, ok := err.(*types.ConditionalCheckFailedException); ok {
			return fmt.Errorf("user stats already exist")
		}
		return fmt.Errorf("failed to create user stats: %w", err)
	}

	return nil
}

// UpdateUserStats updates existing user statistics
func (db *DB) UpdateUserStats(ctx context.Context, stats *UserStats) error {
	// Marshal the stats to DynamoDB format
	item, err := attributevalue.MarshalMap(stats)
	if err != nil {
		return fmt.Errorf("failed to marshal user stats: %w", err)
	}

	_, err = db.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(db.userStatsTableName),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to update user stats: %w", err)
	}

	return nil
}
