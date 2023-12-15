package files

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"cekkustomer.com/pkg/aws"
	"github.com/gin-gonic/gin"
)

func ImportExcel(ctx *gin.Context) {

	localUploadDir := "./uploads"
	uploadFolder := "folder-user"
	bucketName := "importxclxit"

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileName := file.Filename

	// check file
	allowedExtension := []string{".xlsx", ".xls"}
	if !isAllowedExtension(fileName, allowedExtension) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "file format only xlsx and xls"})
		return
	}

	filePath := filepath.Join(localUploadDir, fileName)

	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	s3FilePath := filepath.Join(uploadFolder, fileName)

	// upload file to s3
	if err := aws.NewConnect().S3.UploadFile(bucketName, filePath, s3FilePath); err != nil {
		log.Println("Failed to upload file to S3:", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := os.Remove(filePath); err != nil {
		log.Println("Failed to remove uploaded file:", err.Error())
	} else {
		log.Println("File removed successfully:", filePath)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "successfully upload",
	})

}

func isAllowedExtension(fileName string, allowedExtension []string) bool {
	for _, ext := range allowedExtension {
		if strings.HasSuffix(fileName, ext) {
			return true
		}
	}

	return false
}
