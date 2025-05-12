package session

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DynamoStore implements Store using DynamoDB, with encryption at rest.
type DynamoStore struct {
	TableName string
	Client    *dynamodb.Client
}

func NewDynamoStore(ctx context.Context) (*DynamoStore, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	client := dynamodb.NewFromConfig(cfg)
	table := os.Getenv("SESSION_TABLE")
	return &DynamoStore{TableName: table, Client: client}, nil
}

func (d *DynamoStore) Put(ctx context.Context, s *Session) error {
	item, err := attributevalue.MarshalMap(s)
	if err != nil {
		return err
	}
	_, err = d.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(d.TableName),
		Item:      item,
	})
	return err
}

func (d *DynamoStore) Get(ctx context.Context, sessionID string) (*Session, error) {
	out, err := d.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(d.TableName),
		Key: map[string]ddbtypes.AttributeValue{
			"session_id": &ddbtypes.AttributeValueMemberS{Value: sessionID},
		},
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, nil
	}
	var s Session
	err = attributevalue.UnmarshalMap(out.Item, &s)
	return &s, err
}

func (d *DynamoStore) Update(ctx context.Context, s *Session) error {
	return d.Put(ctx, s)
}
