package service

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
	"github.com/benkim0414/geoauth2/client"
	"github.com/benkim0414/geoauth2/internal/rand"
	"github.com/go-kit/kit/log"
	"github.com/google/uuid"
)

// Service is an interface representing the ability to execute an authorization
// process, obtaining the access token for a given authorization grant.
type Service interface {
	PostClient(ctx context.Context, c *client.Client) (*client.Client, error)
	GetClient(ctx context.Context, id string) (*client.Client, error)
	DeleteClient(ctx context.Context, id string) error
}

type service struct {
	ddb *dynamodb.DynamoDB
}

// New returns a Service including DynamoDB configured by the given config.
func New(config aws.Config, logger log.Logger) Service {
	ddb := dynamodb.New(config)
	err := client.CreateTable(ddb)
	if err != nil {
		panic(err.Error())
	}

	var svc Service
	svc = &service{
		ddb: ddb,
	}
	svc = LoggingMiddleware(logger)(svc)
	return svc
}

func (s *service) PostClient(_ context.Context, c *client.Client) (*client.Client, error) {
	c.ID = uuid.New().String()
	secret, err := rand.HexEncodedString(32)
	if err != nil {
		return nil, err
	}
	c.Secret = secret
	av, err := dynamodbattribute.MarshalMap(c)
	if err != nil {
		return nil, err
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(client.TableName),
	}
	req := s.ddb.PutItemRequest(input)
	_, err = req.Send()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (s *service) GetClient(_ context.Context, id string) (*client.Client, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(client.TableName),
	}
	req := s.ddb.GetItemRequest(input)
	result, err := req.Send()
	if err != nil {
		return nil, err
	}
	c := &client.Client{}
	err = dynamodbattribute.UnmarshalMap(result.Item, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (s *service) DeleteClient(_ context.Context, id string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(client.TableName),
	}
	req := s.ddb.DeleteItemRequest(input)
	_, err := req.Send()
	return err
}
