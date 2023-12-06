package database

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type AwsDynamoDB struct {
	client *dynamodb.Client
}

func NewDynamoDBConnect(config *aws.Config) *AwsDynamoDB {
	client := dynamodb.NewFromConfig(*config)
	return &AwsDynamoDB{
		client: client,
	}
}
