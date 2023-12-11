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

	// check file
	allowedExtension := []string{".xlsx", ".xls"}
	if !isAllowedExtension(fileName, allowedExtension) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "file format only xlsx and xls"})
		return
	}

	filePath := filepath.Join(uploadFolder, fileName)

	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// check file exists
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

		ctx.JSON(http.StatusOK, gin.H{
			"location_file": uploadFile,
		})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "nama file telah ada"})
	}

}

func isAllowedExtension(fileName string, allowedExtension []string) bool {
	for _, ext := range allowedExtension {
		if strings.HasSuffix(fileName, ext) {
			return true
		}
	}

	return false
}
