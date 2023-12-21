package auth

import (
	"fmt"
	"net/http"

	"cekkustomer.com/dtos"
	"cekkustomer.com/pkg/aws"
	"github.com/gin-gonic/gin"
)

func Register(ctx *gin.Context) {
	var registerData dtos.AuthData
	if err := ctx.ShouldBindJSON(&registerData); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	result, err := aws.NewConnect().Cognito.SignUp(registerData.Username, registerData.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("kode konfirmasi telah kami kirim ke %s", result),
	})
}
