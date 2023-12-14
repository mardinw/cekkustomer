package files

import (
	"fmt"
	"net/http"

	"cekkustomer.com/pkg/aws"
	"github.com/gin-gonic/gin"
)

func DeleteFile(ctx *gin.Context) {
	bucketName := "importxclxit"
	fileName := ctx.Param("filename")
	folderUser := ctx.Param("foldername")

	filePath := fmt.Sprintf("%s/%s", folderUser, fileName)

	if err := aws.NewConnect().S3.DeleteFile(bucketName, filePath); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("delete file %s successfully", filePath)})
}
