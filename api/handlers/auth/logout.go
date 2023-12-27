package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"cekkustomer.com/configs"
	"cekkustomer.com/pkg/aws"
	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-envconfig"
)

func Logout(ctx *gin.Context) {
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

	if err := aws.NewConnect().DynamoDB.DeleteTTLSession(tableDynamo.TTLSes, key); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete token from table"})
		return
	}

	expiredCookie := http.Cookie{
		Name:     "access_token",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour),
	}

	http.SetCookie(ctx.Writer, &expiredCookie)
	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
