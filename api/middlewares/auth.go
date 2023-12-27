package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"cekkustomer.com/configs"
	"cekkustomer.com/pkg/aws"
	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-envconfig"
)

func Auth(ctx *gin.Context) {
	var tableDynamo configs.AwsDynTblConfig

	if err := envconfig.Process(ctx, &tableDynamo); err != nil {
		log.Fatal(err.Error())
	}

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
	output, err := aws.NewConnect().Cognito.GetUsername(accessToken)
	if err != nil {
		log.Println(err.Error())
		return
	}

	key := fmt.Sprintf(*output.Username)
	cachedToken, err := aws.NewConnect().DynamoDB.GetSessionByUUID(tableDynamo.TTLSes, key)
	if err != nil || cachedToken.AccessToken != accessToken {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	ctx.Set("uuid", *output.Username)
	ctx.Next()
}
