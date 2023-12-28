package cekdata

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sync"

	"cekkustomer.com/api/models"
	"cekkustomer.com/configs"
	"cekkustomer.com/dtos"
	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-envconfig"
	"golang.org/x/net/context"
)

func GetKec(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		result, err := models.GetAllKec(db)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"locate": result})
	}
}

func CheckDPT(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var cekMatch models.ImportCustomerXls
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

		getKec, err := models.GetAllKec(db)
		if err != nil {
			log.Println(err.Error())
			return
		}

		// readfile excel
		fileName := ctx.Param("filename")
		filePath := fmt.Sprintf("%s/%s", uuidStr, fileName)

		// cek nama di table dpt
		firstName := ctx.Query("nama")
		agenciesName := uuidStr
		// results := make(map[string]interface{})

		// buat buffered channelnya
		resultChannel := make(chan map[string]interface{}, len(getKec))

		var wg sync.WaitGroup

		for _, tableName := range getKec {
			wg.Add(1)

			go func(tableName string) {
				defer wg.Done()

				var err error
				var result []dtos.CheckDPT

				if firstName == "" {
					result, err = cekMatch.GetAll(db, tableName, agenciesName, filePath)
				} else {
					result, err = cekMatch.GetAllByName(db, tableName, agenciesName, firstName, filePath)
				}
				if err != nil {
					ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
					return
				}

				if len(result) > 0 {
					resultChannel <- map[string]interface{}{tableName: result}
				}
			}(tableName)
		}

		go func() {
			wg.Wait()
			close(resultChannel)
		}()

		results := make(map[string]interface{})
		for res := range resultChannel {
			for key, value := range res {
				results[key] = value
			}
		}

		ctx.IndentedJSON(http.StatusOK, gin.H{
			"results": results,
		})

	}
}
