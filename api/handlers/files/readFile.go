package files

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"cekkustomer.com/api/middlewares"
	"github.com/gin-gonic/gin"
)

func ReadFile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileName := file.Filename
	filePath := filepath.Join(uploadFolder, fileName)
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

	ctx.IndentedJSON(http.StatusOK, gin.H{"data": jsonData})
}
