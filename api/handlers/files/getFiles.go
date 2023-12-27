package files

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"cekkustomer.com/configs"
	"cekkustomer.com/pkg/aws"
	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-envconfig"
	"golang.org/x/net/context"
)

func DownloadSampleXlsx(ctx *gin.Context) {
	var bucketFolder configs.AwsS3Bucket
	if err := envconfig.Process(context.Background(), &bucketFolder); err != nil {
		log.Fatal(err.Error())
	}

	fileName := ctx.Param("filename")

	s3FilePath := "sample/" + fileName

	if !aws.NewConnect().S3.CheckExists(ctx, bucketFolder.ImportS3, s3FilePath) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "file not found in s3."})
		return
	}

	err := aws.NewConnect().S3.DownloadFile(bucketFolder.ImportS3, s3FilePath, fileName)
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

	var bucketFolder configs.AwsS3Bucket
	if err := envconfig.Process(context.Background(), &bucketFolder); err != nil {
		log.Fatal(err.Error())
	}

	uuid, exists := ctx.Get("uuid")
	if !exists {
		log.Println("uuid tidak ditemukan")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	uuidStr, ok := uuid.(string)
	if !ok {
		log.Println("gagal konversi ke string")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	objects, err := aws.NewConnect().S3.ListFile(bucketFolder.ImportS3, uuidStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Failed to list objects: %v", err.Error())})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"files": objects})
}
