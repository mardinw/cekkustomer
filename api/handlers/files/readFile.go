package files

import (
	"net/http"

	"cekkustomer.com/api/middlewares"
	"cekkustomer.com/pkg/aws"
	"github.com/gin-gonic/gin"
)

func ReadFile(ctx *gin.Context) {
	bucketName := "importxclxit"
	fileName := "sample_cekkustomer.xlsx"

	getFile, err := aws.NewConnect().S3.GetFile(bucketName, fileName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	readFile, err := middlewares.ReadExcel(getFile.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": readFile})
}