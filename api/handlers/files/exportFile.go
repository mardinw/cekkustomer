package files

import (
	"net/http"

	"cekkustomer.com/api/middlewares"
	"github.com/gin-gonic/gin"
)

func ExportToExcel(ctx *gin.Context) {
	jsonData, exists := ctx.Get("jsonData")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "json data not found"})
		return
	}

	if err := middlewares.CreateExcel(jsonData.(string)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create excel file"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "file create successfully"})
}
