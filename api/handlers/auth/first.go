package auth

import (
	"net/http"

	"cekkustomer.com/pkg/aws"
	"github.com/gin-gonic/gin"
)

func FirstLogin(ctx *gin.Context) {
	username := ctx.Param("username")

	if err := aws.NewConnect().Cognito.PreferredUsernameAttributes("test saja", username); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}
	ctx.Next()
}
