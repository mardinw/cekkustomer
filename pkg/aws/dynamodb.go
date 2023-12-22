package aws

import (
	"context"
	"errors"
	"fmt"
	"log"

	"cekkustomer.com/dtos"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type AwsDynamoDB struct {
	DynamoDBClient *dynamodb.Client
	TableName      string
}

func NewDynamoDBConnect(config *aws.Config) *AwsDynamoDB {
	client := dynamodb.NewFromConfig(*config)
	return &AwsDynamoDB{
		DynamoDBClient: client,
	}
}

func (c *AwsDynamoDB) TableExists(tableName string) (bool, error) {
	exists := true
	_, err := c.DynamoDBClient.DescribeTable(
		context.TODO(),
		&dynamodb.DescribeTableInput{
			TableName: aws.String(tableName),
		},
	)
	if err != nil {
		var notFoundEx *types.ResourceNotFoundException
		if errors.As(err, &notFoundEx) {
			log.Printf("Table %v does not exists.\n", tableName)
			err = nil

		} else {
			log.Printf("Couldn't determine existence of table %v. Here's why: %v\n", tableName, err)
		}
		exists = false
	}

	return exists, err
}

func (c *AwsDynamoDB) SaveTTLSession(tableName, uuid, token string, createdAt, expireAt int64) error {
	item := map[string]types.AttributeValue{
		"uuid": &types.AttributeValueMemberS{
			Value: uuid,
		},
		"access_token": &types.AttributeValueMemberS{
			Value: token,
		},
		"created_at": &types.AttributeValueMemberN{
			Value: fmt.Sprintf("%d", createdAt),
		},
		"expire_at": &types.AttributeValueMemberN{
			Value: fmt.Sprintf("%d", expireAt),
		},
	}

	_, err := c.DynamoDBClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	if err != nil {
		log.Printf("couldn't add item to table. Here's why: %v\n", err)
	}
	return err
}

func (c *AwsDynamoDB) GetSessionByUUID(tableName, uuid string) (*dtos.TTLSessionData, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"uuid": &types.AttributeValueMemberS{
				Value: uuid,
			},
		},
	}

	result, err := c.DynamoDBClient.GetItem(context.TODO(), input)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if result.Item == nil {
		return nil, fmt.Errorf("item with UUID %s not found", uuid)
	}

	session := &dtos.TTLSessionData{}
	err = attributevalue.UnmarshalMap(result.Item, session)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return session, nil
}

func (c *AwsDynamoDB) DeleteTTLSession(tableName, uuid string) error {

	key := map[string]types.AttributeValue{
		"uuid": &types.AttributeValueMemberS{
			Value: uuid,
		},
	}

	_, err := c.DynamoDBClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key:       key,
	})
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
