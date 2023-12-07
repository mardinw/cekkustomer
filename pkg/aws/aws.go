package aws

import (
	"context"
	"log"

	"cekkustomer.com/configs"
	"cekkustomer.com/pkg/database"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/sethvargo/go-envconfig"
)

type AwsConnect struct {
	S3       *AwsS3
	DynamoDB *database.AwsDynamoDB
}

func NewConnect() *AwsConnect {
	ctx := context.Background()

	var configs configs.AppConfiguration

	if err := envconfig.Process(ctx, &configs); err != nil {
		log.Fatal(err.Error())
	}

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(configs.AwsConf.AwsRegion),
		config.WithSharedConfigProfile(configs.AwsConf.AwsProfile),
	)

	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return &AwsConnect{
		S3:       NewS3Connect(&cfg),
		DynamoDB: database.NewDynamoDBConnect(&cfg),
	}
}
