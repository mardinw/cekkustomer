package auth

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"cekkustomer.com/configs"
	"cekkustomer.com/dtos"
	"cekkustomer.com/pkg/aws"
	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-envconfig"
)

func Login(ctx *gin.Context) {
	var loginData dtos.AuthData
	var bucketS3 configs.AwsS3Bucket
	var tableDynamo configs.AwsDynTblConfig

	if err := envconfig.Process(ctx, &bucketS3); err != nil {
		log.Fatal(err.Error())
	}

	if err := envconfig.Process(ctx, &tableDynamo); err != nil {
		log.Fatal(err.Error())
	}

	if err := ctx.ShouldBindJSON(&loginData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := aws.NewConnect().Cognito.SignIn(loginData.Username, loginData.Password)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err.Error())
		return
	}

	output, err := aws.NewConnect().Cognito.GetUsername(*result.AccessToken)
	if err != nil {
		log.Println(err.Error())
		return
	}
	folderName := fmt.Sprintf("%s/", *output.Username)

	isExists, err := aws.NewConnect().S3.CheckFolderExistsInBucket(bucketS3.ImportS3, folderName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if !isExists {
		if err := aws.NewConnect().S3.CreateFolderInBucket(bucketS3.ImportS3, folderName); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// create time create session and expire
	createdAt := time.Now().UnixMilli()
	expireAt := time.Now().Add(time.Hour).UnixNano() / int64(time.Millisecond)

	// save session
	if err := aws.NewConnect().DynamoDB.SaveTTLSession(tableDynamo.TTLSes, *output.Username, *result.AccessToken, createdAt, expireAt); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	session := http.Cookie{
		Name:     "access_token",
		Value:    *result.AccessToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 1),
	}

	http.SetCookie(ctx.Writer, &session)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "successfully login",
	})

}
