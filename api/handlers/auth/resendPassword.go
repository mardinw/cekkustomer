package auth

import (
	"net/http"

	"cekkustomer.com/dtos"
	"cekkustomer.com/pkg/aws"
	"github.com/gin-gonic/gin"
)

func ResetPassword(ctx *gin.Context) {
	var resetPassword dtos.AuthResetData
	if err := ctx.BindJSON(&resetPassword); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	_, err := aws.NewConnect().Cognito.ResetPassword(resetPassword.Username, resetPassword.Password, resetPassword.ConfirmationCode)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "reset password telah sukses"})
}
