package files

import (
	"fmt"
	"net/http"

	"cekkustomer.com/pkg/aws"
	"github.com/gin-gonic/gin"
)

func DownloadSampleXlsx(ctx *gin.Context) {
	fileName := ctx.Param("filename")

	bucketName := "importxclxit"

	downFile := aws.NewConnect().S3.DownloadFile(bucketName, fileName)
	if downFile != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": downFile.Error()})
		return
	}

	ctx.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	ctx.Writer.Header().Set("Content-Type", "application/octet-stream")
	ctx.File(fileName)
}
