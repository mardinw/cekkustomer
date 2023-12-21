package auth

import (
	"fmt"
	"net/http"

	"cekkustomer.com/dtos"
	"cekkustomer.com/pkg/aws"
	"github.com/gin-gonic/gin"
)

func ResendCode(ctx *gin.Context) {
	var resendCode dtos.Users

	if err := ctx.BindJSON(&resendCode); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := aws.NewConnect().Cognito.ResendConfirmationCode(resendCode.Username)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("kode konfirmasi telah dikirim ulang kembali ke %s", result)})
}
