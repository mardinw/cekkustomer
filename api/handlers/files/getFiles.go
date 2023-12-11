package files

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"cekkustomer.com/pkg/aws"
	"github.com/gin-gonic/gin"
)

func DownloadSampleXlsx(ctx *gin.Context) {
	fileName := ctx.Param("filename")

	s3FilePath := "sample/" + fileName

	bucketName := "importxclxit"

	if !aws.NewConnect().S3.CheckExists(ctx, bucketName, s3FilePath) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "file not found in s3."})
		return
	}

	err := aws.NewConnect().S3.DownloadFile(bucketName, s3FilePath, fileName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func() {
		err := os.Remove(fileName)
		if err != nil {
			log.Println(err.Error())
		}
	}()

	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))

	ctx.File(fileName)
}

func GetListFolder(ctx *gin.Context) {
	folder := ctx.Param("folder")
	bucketName := "importxclxit"

	objects, err := aws.NewConnect().S3.ListFile(bucketName, folder)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Failed to list objects: %v", err.Error())})
		return
	}

	ctx.JSON(http.StatusOK, objects)
}
