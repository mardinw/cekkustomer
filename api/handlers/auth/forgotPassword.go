package auth

import (
	"fmt"
	"net/http"

	"cekkustomer.com/dtos"
	"cekkustomer.com/pkg/aws"
	"github.com/gin-gonic/gin"
)

func ForgotPassword(ctx *gin.Context) {
	var forgotData dtos.Users

	if err := ctx.BindJSON(&forgotData); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	result, err := aws.NewConnect().Cognito.ForgotPassword(forgotData.Username)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("kode konfirmasi untuk lupa password telah dikirim, silahkan cek inbox %s", result)})
}
