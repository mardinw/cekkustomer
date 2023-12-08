package files

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"cekkustomer.com/pkg/aws"
	"github.com/gin-gonic/gin"
)

const (
	uploadFolder = "external/importxclxit"
)

func ImportExcel(ctx *gin.Context) {

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileName := file.Filename
	filePath := filepath.Join(uploadFolder, fileName)

	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	checkFile := aws.NewConnect().S3.CheckExists(ctx, "importxclxit", fileName)

	if !checkFile {
		uploadFile, err := aws.NewConnect().S3.UploadFile(ctx, "importxclxit", fileName, filePath)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := os.Remove(filePath); err != nil {
			log.Println("Failed to remove uploaded file:", err.Error())
		} else {
			log.Println("File removed successfully:", filePath)
		}

		uploadedTime := time.Now().UnixMicro()
		agencies := "mardin"

		// record to dynamodb
		err = aws.NewConnect().DynamoDB.AddImportXlsx("import_xlsx", agencies, uploadFile, uploadedTime)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"location_file": uploadFile,
		})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "nama file telah ada"})
	}

}
