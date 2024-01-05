package files

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"cekkustomer.com/api/middlewares"
	"cekkustomer.com/api/models"
	"cekkustomer.com/configs"
	"cekkustomer.com/pkg/aws"
	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-envconfig"
	"golang.org/x/net/context"
)

func isAllowedExtension(fileName string, allowedExtension []string) bool {
	for _, ext := range allowedExtension {
		if strings.HasSuffix(fileName, ext) {
			return true
		}
	}

	return false
}

func ImportExcel(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {

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

		localUploadDir := "./uploads"
		agenciesName := uuidStr

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
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		s3FilePath := filepath.Join(uuidStr, fileName)

		// upload file to s3
		if err := aws.NewConnect().S3.UploadFile(bucketFolder.ImportS3, filePath, s3FilePath); err != nil {
			log.Println("Failed to upload file to S3:", err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := os.Remove(filePath); err != nil {
			log.Println("Failed to remove uploaded file:", err.Error())
		} else {
			log.Println("File removed successfully:", filePath)
		}

		// read file from bucket
		getFile, err := aws.NewConnect().S3.GetFile(bucketFolder.ImportS3, s3FilePath)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		dataChannel := make(chan map[string]interface{}, 10)

		readFile, err := middlewares.ReadExcel(getFile.Body, bucketFolder.ImportS3, s3FilePath)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// end readfile
		// start insert to table customer

		// goroutine for import data concurrently
		go func() {
			for _, data := range readFile {
				dataChannel <- data
			}

			close(dataChannel)
		}()

		// Goroutines to process data concurrently
		var wg sync.WaitGroup
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				processData(db, agenciesName, s3FilePath, dataChannel)
			}()
		}

		wg.Wait()

		ctx.JSON(http.StatusOK, gin.H{
			"message": "successfully uploaded",
		})
	}
}

func processData(db *sql.DB, agenciesName, s3FilePath string, dataChannel <-chan map[string]interface{}) {
	for data := range dataChannel {
		timeNow := time.Now().UnixMilli()

		// card number
		cardNumberValue := data["card_number"]
		var cardNumber string
		if cardNumberValue != nil {
			cardNumberStr, ok := data["card_number"].(string)
			if !ok || cardNumber == "" {
				log.Println("card number not found")
				cardNumber = ""
			}
			cardNumber = cardNumberStr
		} else {
			cardNumber = ""
		}

		// nik number
		nikNumberValue := data["nik"]
		var nikNumber string
		var nikHidden string
		// gunakan regex untuk compile
		re := regexp.MustCompile("[0-9]+")

		if nikNumberValue != nil {
			nikNumberRaw, ok := data["nik"].(string)

			if !ok || nikNumberRaw == "" {
				log.Println("nik tidak ada")
				nikNumber = ""
			}
			matches := re.FindAllString(nikNumberRaw, -1)
			// ambil angka
			if len(matches) > 0 {
				nikNumberMatch := matches[0]

				// cek length
				if len(nikNumberMatch) > 4 {
					// ambil substring dari string
					hiddenPart := strings.Repeat("*", 4)
					visiblePart := nikNumberMatch[:len(nikNumberMatch)-4]

					// Gabungkan bagian yang terlihat dan tersembunyi
					nikHidden = visiblePart + hiddenPart
					log.Println(nikHidden)
				}
				nikNumber = nikNumberMatch
			}

		}

		// collector
		collectorRaw := data["collector"]
		var collector string
		if collectorRaw != nil {
			collectorStr, ok := collectorRaw.(string)
			if !ok {
				log.Println("collector is not a string")
				collector = ""
			}
			collector = collectorStr
		} else {
			collector = ""
		}

		// cek concat customer
		concatCustomer := data["concat_customer"]
		var concatCustToUpper string

		if concatCustomer != nil {
			concatCustomerValue, ok := data["concat_customer"].(string)
			if !ok || concatCustomerValue == "" {
				log.Println("concat customer not found")
				concatCustomerValue = ""
			}
			concatCustToUpper = strings.ToUpper(concatCustomerValue)
		} else {
			concatCustToUpper = ""
		}

		firstName := data["first_name"]
		var firstNameValue string
		if firstName != nil {
			firstNameStr, ok := data["first_name"].(string)
			if !ok || firstNameValue == "" {
				log.Println("first name not found")
				firstNameValue = ""
			}
			firstNameValue = firstNameStr
		} else {
			firstNameValue = ""
		}

		address3 := data["address_3"]
		var address3Value string
		if address3 != nil {
			address3Str, ok := data["address_3"].(string)
			if !ok || address3Value == "" {
				log.Println("address 3 not found")
				address3Value = ""
			}
			address3Value = address3Str
		} else {
			address3Value = ""
		}

		address4 := data["address_4"]
		var address4Value string
		if address4 != nil {
			address4Str, ok := data["address_4"].(string)
			if !ok || address4Value == "" {
				log.Println("address 3 not found")
				address4Value = ""
			}
			address4Value = address4Str
		} else {
			address4Value = ""
		}

		zipCode := data["zipcode"]
		var zipCodeValue string
		if zipCode != nil {
			zipCodeStr, ok := data["zipcode"].(string)
			if !ok || zipCodeValue == "" {
				log.Println("zipcode not found")
				zipCodeValue = ""
			}
			zipCodeValue = zipCodeStr
		} else {
			zipCodeValue = ""
		}

		inputCustomer := &models.ImportCustomerXls{
			CardNumber:     cardNumber,
			NIK:            nikNumber,
			NIKCheck:       nikHidden,
			FirstName:      firstNameValue,
			Collector:      collector,
			Agencies:       agenciesName,
			Address3:       address3Value,
			Address4:       address4Value,
			ZipCode:        zipCodeValue,
			ConcatCustomer: concatCustToUpper,
			Files:          s3FilePath,
			Created:        timeNow,
		}

		if err := inputCustomer.InsertCustomer(db); err != nil {
			log.Println("failed to insert customer:", err.Error())
			return
		}

	}
}
