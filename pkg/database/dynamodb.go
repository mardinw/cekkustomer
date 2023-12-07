package database

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
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

func (c *AwsDynamoDB) AddImportXlsx(tableName string, agencies, uploadFile string, uploadedTime int64) error {

	item := map[string]types.AttributeValue{
		"agencies": &types.AttributeValueMemberS{
			Value: agencies,
		},
		"files":    &types.AttributeValueMemberS{Value: uploadFile},
		"uploaded": &types.AttributeValueMemberN{Value: fmt.Sprint(uploadedTime)},
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

func (c *AwsDynamoDB) AddCustomerXlsx(tableName, agencies, firstName, addressThree, addressFour, concatCustomer, collector string, cardNumber, created int64, homeZipCode int32) error {

	item := map[string]types.AttributeValue{
		"agencies": &types.AttributeValueMemberS{
			Value: agencies,
		},
		"card_number": &types.AttributeValueMemberN{
			Value: fmt.Sprint(cardNumber),
		},
		"first_name": &types.AttributeValueMemberS{
			Value: firstName,
		},
		"address_3": &types.AttributeValueMemberS{
			Value: addressThree,
		},
		"address_4": &types.AttributeValueMemberS{
			Value: addressFour,
		},
		"home_zip_code": &types.AttributeValueMemberN{
			Value: fmt.Sprint(homeZipCode),
		},
		"concat_customer": &types.AttributeValueMemberS{
			Value: concatCustomer,
		},
		"collector": &types.AttributeValueMemberS{
			Value: collector,
		},
		"created": &types.AttributeValueMemberN{
			Value: fmt.Sprint(created),
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
