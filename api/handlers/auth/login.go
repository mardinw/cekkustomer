package auth

import (
	"log"
	"net/http"

	"cekkustomer.com/dtos"
	"cekkustomer.com/pkg/aws"
	"github.com/gin-gonic/gin"
)

func Login(ctx *gin.Context) {
	var loginData dtos.AuthData
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

	ctx.JSON(http.StatusOK, gin.H{
		"username": *output.Username,
	})
}
