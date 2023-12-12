package files

import (
	"fmt"
	"net/http"

	"cekkustomer.com/api/middlewares"
	"cekkustomer.com/pkg/aws"
	"github.com/gin-gonic/gin"
)

func ReadFile(ctx *gin.Context) {
	bucketName := "importxclxit"
	fileName := ctx.Param("filename")
	folderUser := ctx.Param("foldername")

	filePath := fmt.Sprintf("%s/%s", folderUser, fileName)

	getFile, err := aws.NewConnect().S3.GetFile(bucketName, filePath)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	readFile, err := middlewares.ReadExcel(getFile.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, readFile)
}
