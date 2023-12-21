package auth

import (
	"net/http"

	"cekkustomer.com/dtos"
	"cekkustomer.com/pkg/aws"
	"github.com/gin-gonic/gin"
)

func Confirmation(ctx *gin.Context) {
	var confirmData dtos.AuthCodeData

	if err := ctx.BindJSON(&confirmData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := aws.NewConnect().Cognito.ConfirmSignUp(confirmData.Username, confirmData.ConfirmationCode)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "akun telah dikonfirmasi"})
}
