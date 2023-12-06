package files

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"cekkustomer.com/api/middlewares"
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

	readFile, err := middlewares.ReadExcel(filePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	// convert the result json
	jsonData, err := json.Marshal(readFile)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read json"})
		return
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
		ctx.JSON(http.StatusOK, gin.H{
			"location_file": uploadFile,
			"data":          jsonData,
		})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "nama file telah ada"})
	}

}
