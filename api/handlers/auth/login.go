package auth

import (
	"fmt"
	"log"
	"net/http"

	"cekkustomer.com/configs"
	"cekkustomer.com/dtos"
	"cekkustomer.com/pkg/aws"
	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-envconfig"
)

func Login(ctx *gin.Context) {
	var loginData dtos.AuthData
	var config configs.AwsS3Bucket
	if err := envconfig.Process(ctx, &config); err != nil {
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

	isExists, err := aws.NewConnect().S3.CheckFolderExistsInBucket(config.ImportS3, folderName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if !isExists {
		if err := aws.NewConnect().S3.CreateFolderInBucket(config.ImportS3, folderName); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "successfully login",
	})

}
