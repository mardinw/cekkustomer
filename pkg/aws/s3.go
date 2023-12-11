package aws

import (
	"context"
	"errors"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/gin-gonic/gin"
)

type AwsS3 struct {
	client *s3.Client
}

func NewS3Connect(config *aws.Config) *AwsS3 {
	client := s3.NewFromConfig(*config)
	return &AwsS3{
		client: client,
	}
}

var apiError smithy.APIError

func (c *AwsS3) DownloadFile(bucketName, objectKey, fileName string) error {

	resp, err := c.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})

	if err != nil {
		log.Println(err.Error())
		return err
	}

	defer resp.Body.Close()

	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("couldn't create file %v. Here why : %v\n", file, err.Error())
		return err
	}

	defer file.Close()

	fileContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	_, err = file.Write(fileContent)
	return err
}

func (c *AwsS3) ListFile(bucketName, folder string) ([]types.Object, error) {
	result, err := c.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(folder),
	})

	var contents []types.Object
	if err != nil {
		log.Printf("Couldn't list objects in bucket %v. Here's why: %v\n", bucketName, err)
	} else {
		contents = result.Contents
	}

	return contents, err
}

func (c *AwsS3) GetFile(bucketName, fileName string) (*s3.GetObjectOutput, error) {
	var err error

	result, err := c.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	})

	defer func() {
		if err != nil {
			result.Body.Close()
		}
	}()

	if err != nil {
		log.Printf("Couldn't get object %v:%v. Here's why: %v\n", bucketName, fileName, err)
		return nil, err
	}

	return result, err
}

func (c *AwsS3) UploadFile(ctx *gin.Context, bucketName, fileName, filePath string) (string, error) {

	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err.Error())
	}
	defer file.Close()

	uploader := manager.NewUploader(c.client)

	res, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Body:   file,
		Key:    aws.String(fileName),
	})

	if err != nil {
		log.Println(err)
		return "", err
	}

	return res.Location, err
}

func (c *AwsS3) CheckExists(ctx *gin.Context, bucketName, fileName string) bool {
	_, err := c.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	})

	exists := true

	if err != nil {
		if errors.As(err, &apiError) {
			switch apiError.(type) {
			case *types.NotFound:
				log.Printf("Bucket %v is available", bucketName)
				exists = false
				err = nil
			default:
				log.Printf("Either you don't have access to bucket %v or another error occured."+
					"Here's what happened: %v\n", bucketName, err)
			}
		}
	} else {
		log.Printf("File %v exists and you already own it.", fileName)
	}

	return exists
}

func (c *AwsS3) ListBuckets(ctx *gin.Context) (*s3.ListBucketsOutput, error) {
	return c.client.ListBuckets(ctx, &s3.ListBucketsInput{})
}

func (c *AwsS3) CreateBucket(ctx *gin.Context, bucketName string) error {
	_, err := c.client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		return err
	}

	return nil
}

func (c *AwsS3) BucketExists(ctx *gin.Context, bucketName string) (bool, error) {
	_, err := c.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})

	exists := true
	if err != nil {
		if errors.As(err, &apiError) {
			switch apiError.(type) {
			case *types.NotFound:
				log.Printf("Bucket %v is available", bucketName)
				exists = false
				err = nil
			default:
				log.Printf("Either you don't have access to bucket %v or another error occured."+
					"Here's what happened: %v\n", bucketName, err)
			}
		}
	} else {
		log.Printf("Bucket %v exists and you already own it.", bucketName)
	}

	return exists, err
}

func (c *AwsS3) DeleteFile(ctx *gin.Context, bucketName, fileName string) error {
	_, err := c.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	})

	if err != nil {
		return err
	}

	return nil
}
