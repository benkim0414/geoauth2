package service

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
	"github.com/go-kit/kit/log"
)

// GrantType is the type of authorization being used by the client to obtain
// an access token.
const (
	// GrantTypeAuthorizationCode is used by confidential and public clients
	// to exchange an authorization code for an access token.
	GrantTypeAuthorizationCode = "authorization_code"

	// GrantTypeImplicit is a simplified flow that can be used by public
	// clients, where the access token is returned immediately without an
	// extra authorization code exchange step.
	GrantTypeImplicit = "implicit"

	// GrantTypePassword is used by first-party clients to exchange a user's
	// credentials for an access token.
	GrantTypePassword = "password"

	// GrantTypeClientCredentials is used by clients to obtain an access
	// token outside of the context of a user.
	GrantTypeClientCredentials = "client_credentials"

	// GrantTypeRefreshToken is used by clients to exchange a refresh token
	// for an access token when the access token has expired.
	GrantTypeRefreshToken = "refresh_token"
)

// ResponseType is the type of response being used by the authorization code
// grant type and implicit grant type flows.
const (
	// ResponseTypeCode is used for requesting an authorization code.
	ResponseTypeCode = "code"
	// ResponseTypeToken is used for requesting an access token.
	ResponseTypeToken = "token"
)

// Client represents an OAuth 2.0 client.
type Client struct {
	// ID is the identifier for this client.
	ID string `json:"id"`

	// Name is the human-readable string name of the client to be presented to the
	// end-user during authorization.
	Name string `json:"name"`

	// Secret is the client's secret.
	Secret string `json:"secret"`

	// RedirectURI is an allowed redirect url for the client.
	RedirectURI string `json:"redirectUri"`

	// GrantType is grant type the client is allowed to use.
	GrantType string `json:"grantType"`

	// ResponseType is the OAuth 2.0 response type string that the client can use at
	// the authorization endpoint.
	ResponseType string `json:"responseType"`

	// Scope is a string containing a space-separated list of scope values (as
	// described in Section 3.3 of OAuth 2.0 [RFC6749]) that the client can use when
	// requesting access tokens.
	Scope string `json:"scope"`

	// Public is a boolean that identifies this client as public, meaning that it
	// does not have a secret. It will disable the client_credentials grant type for
	// this client if set.
	Public bool `json:"public"`
}

// Service is an interface representing the ability to execute an authorization
// process, obtaining the access token for a given authorization grant.
type Service interface {
	PostClient(ctx context.Context, c *Client) (*Client, error)
	GetClient(ctx context.Context, id string) (*Client, error)
	DeleteClient(ctx context.Context, id string) error
}

type service struct {
	ddb *dynamodb.DynamoDB
}

// New returns a Service including DynamoDB configured by the given config.
func New(config aws.Config, logger log.Logger) Service {
	ddb := dynamodb.New(config)
	err := createClientTable(ddb)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceInUseException:
				logger.Log("dynamodb: %v", aerr.Error())
			default:
				panic(aerr.Error())
			}
		} else {
			panic(err.Error())
		}
	}

	var svc Service
	svc = &service{
		ddb: ddb,
	}
	svc = LoggingMiddleware(logger)(svc)
	return svc
}

const (
	TableNameClients = "clients"
)

func createClientTable(ddb *dynamodb.DynamoDB) error {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: dynamodb.ScalarAttributeTypeS,
			},
		},
		KeySchema: []dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       dynamodb.KeyTypeHash,
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: aws.String(TableNameClients),
	}
	req := ddb.CreateTableRequest(input)
	_, err := req.Send()
	return err
}

func (s *service) PostClient(_ context.Context, c *Client) (*Client, error) {
	av, err := dynamodbattribute.MarshalMap(c)
	if err != nil {
		return nil, err
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(TableNameClients),
	}
	req := s.ddb.PutItemRequest(input)
	_, err = req.Send()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (s *service) GetClient(_ context.Context, id string) (*Client, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(TableNameClients),
	}
	req := s.ddb.GetItemRequest(input)
	result, err := req.Send()
	if err != nil {
		return nil, err
	}
	c := &Client{}
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
		TableName: aws.String(TableNameClients),
	}
	req := s.ddb.DeleteItemRequest(input)
	_, err := req.Send()
	return err
}
