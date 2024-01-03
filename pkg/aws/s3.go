package aws

import (
	"context"
	"errors"
	"io"
	"log"
	"os"
	"strings"

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

func (c *AwsS3) ListFile(bucketName, folder string) ([]string, error) {
	result, err := c.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(folder),
	})

	var fileList []string
	// if err != nil {
	// 	log.Printf("Couldn't list objects in bucket %v. Here's why: %v\n", bucketName, err)
	// } else {
	// 	contents = result.Contents
	// }
	for _, item := range result.Contents {
		if !strings.HasSuffix(*item.Key, "/") {
			fileList = append(fileList, *item.Key)
		}
	}

	return fileList, err
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

func (c *AwsS3) UploadFile(bucketName, localFilePath, objectKey string) error {

	file, err := os.Open(localFilePath)
	if err != nil {
		log.Println(err.Error())
	}
	defer file.Close()

	uploader := manager.NewUploader(c.client)

	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Body:   file,
		Key:    aws.String(objectKey),
	})

	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
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
	}

	return exists
}

func (c *AwsS3) ListBuckets(ctx *gin.Context) (*s3.ListBucketsOutput, error) {
	return c.client.ListBuckets(ctx, &s3.ListBucketsInput{})
}

func (c *AwsS3) CreateBucket(bucketName string) error {
	_, err := c.client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		return err
	}

	return nil
}

func (c *AwsS3) BucketExists(bucketName string) (bool, error) {
	_, err := c.client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
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

func (c *AwsS3) DeleteFile(bucketName, objectKey string) error {
	_, err := c.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})

	if err != nil {
		log.Printf("Couldn't delete objects from bucket %v with files %v. Here's why: %v\n", bucketName, objectKey, err)
	}

	return err
}

func (c *AwsS3) CreateFolderInBucket(bucketName, folderName string) error {

	input := &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &folderName,
		Body:   nil,
	}

	_, err := c.client.PutObject(context.TODO(), input)
	if err != nil {
		return err
	}

	return nil
}

func (c *AwsS3) CheckFolderExistsInBucket(bucketName, folderName string) (bool, error) {
	keys := int32(1)
	resp, err := c.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket:  &bucketName,
		Prefix:  &folderName,
		MaxKeys: &keys,
	})
	if err != nil {
		return false, err
	}

	if len(resp.Contents) > 0 && strings.HasPrefix(*resp.Contents[0].Key, folderName) {
		return true, nil
	}

	return false, nil

}
