package cekdata

import (
	"net/http"
	"strings"

	"cekkustomer.com/api/helpers"
	"cekkustomer.com/pkg/aws"
	"github.com/gin-gonic/gin"
)

func GetAttributes(ctx *gin.Context) {

	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	splitted := strings.Split(authHeader, " ")
	if len(splitted) != 2 || strings.ToLower(splitted[0]) != "bearer" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		return
	}

	accessToken := splitted[1]
	outputUser, err := aws.NewConnect().Cognito.GetUsername(accessToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	// check role
	if err := helpers.CheckAccountAdmin(outputUser.Username); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}

	userName := ctx.Param("user")

	output, err := aws.NewConnect().Cognito.CheckUserAttributes(userName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, output)
}
