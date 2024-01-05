package middlewares

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"cekkustomer.com/pkg/aws"
	"github.com/xuri/excelize/v2"
)

type MapCustomer []map[string]interface{}

func ReadExcel(fileName io.ReadCloser, bucketName, s3FilePath string) (MapCustomer, error) {
	xlsx, err := excelize.OpenReader(fileName)
	if err != nil {
		log.Println(err.Error())
	}

	sheetName := xlsx.GetSheetList()[0]

	rows, err := xlsx.GetRows(sheetName)
	if err != nil {
		log.Println(err.Error())
	}

	// cek baris
	if len(rows) > 2000 {
		err := errors.New("oops baris lebih dari 200")

		if err := aws.NewConnect().S3.DeleteFile(bucketName, s3FilePath); err != nil {
			log.Println("file not found")
			return nil, err
		}
		return nil, err
	}
	// konversi data excel to json
	var result MapCustomer
	keys := rows[0]

	checkKeys := []string{
		"card_number",
		"nik",
		"first_name",
		"address_3",
		"address_4",
		"zipcode",
		"collector",
		"concat_customer (nama + tgl lahir)",
	}

	// mengecek table header
	for _, key := range checkKeys {
		switch key {
		case "card_number":
			continue
		case "nik":
			continue
		case "first_name":
			continue
		case "address_3":
			continue
		case "address_4":
			continue
		case "zipcode":
			continue
		case "collector":
			continue
		case "concat_customer":
			continue
		case "concat_customer (nama + tgl lahir)":
			continue
		default:
			err = errors.New("key tidak ditemukan")
			return nil, err
		}
	}

	for _, row := range rows[1:] {
		rowData := make(map[string]interface{})
		for colIndex, cell := range row {
			key := keys[colIndex]
			rowData[key] = cell
		}

		result = append(result, rowData)
	}

	return result, err
}

func CreateExcel(jsonData, bucketExport, fileName, filePath string) error {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return err
	}

	// membuat file excel

	file := excelize.NewFile()

	sheetName := "Match Data"
	headers := []string{"Card Number", "NIK", "First Name", "Collector", "Agencies", "Address 3", "Address 4", "Zip Code", "Kode Pos", "Kelurahan", "Kecamatan", "Nama"}
	file.SetSheetName(file.GetSheetName(0), sheetName)
	file.SetCellValue(sheetName, "A1", "Data Customer")
	file.SetCellValue(sheetName, "I1", "Data Match")
	file.MergeCell(sheetName, "A1", "H1")
	file.MergeCell(sheetName, "I1", "L1")

	// Membuat table header untuk data
	for colIndex, header := range headers {
		colName, _ := excelize.ColumnNumberToName(colIndex + 1)
		file.SetCellValue(sheetName, fmt.Sprintf("%s2", colName), header)
	}

	// tambahkan data ke sheet
	rowIndex := 3
	for _, tableDataArray := range data {
		if tableDataArray == nil {
			continue
		}
		for _, tableData := range tableDataArray.(map[string]interface{}) {
			if tableData == nil {
				continue
			}
			// data customer
			for _, colMapCustomer := range tableData.([]interface{}) {
				for colName, colValueCustomer := range colMapCustomer.(map[string]interface{}) {
					switch colName {
					case "card_number":
						file.SetCellValue(sheetName, fmt.Sprintf("A%d", rowIndex), fmt.Sprintf("'%s", colValueCustomer))
					case "nik":
						file.SetCellValue(sheetName, fmt.Sprintf("B%d", rowIndex), fmt.Sprintf("'%s", colValueCustomer))
					case "first_name":
						file.SetCellValue(sheetName, fmt.Sprintf("C%d", rowIndex), colValueCustomer)
					case "collector":
						file.SetCellValue(sheetName, fmt.Sprintf("D%d", rowIndex), colValueCustomer)
					case "agencies":
						file.SetCellValue(sheetName, fmt.Sprintf("E%d", rowIndex), colValueCustomer)
					case "address_3":
						file.SetCellValue(sheetName, fmt.Sprintf("F%d", rowIndex), colValueCustomer)
					case "address_4":
						file.SetCellValue(sheetName, fmt.Sprintf("G%d", rowIndex), colValueCustomer)
					case "home_zip_code":
						file.SetCellValue(sheetName, fmt.Sprintf("H%d", rowIndex), fmt.Sprintf("'%s", colValueCustomer))
					case "kodepos":
						file.SetCellValue(sheetName, fmt.Sprintf("I%d", rowIndex), fmt.Sprintf("'%s", strconv.FormatFloat(colValueCustomer.(float64), 'f', -1, 64)))
					case "kelurahan":
						file.SetCellValue(sheetName, fmt.Sprintf("J%d", rowIndex), colValueCustomer)
					case "kecamatan":
						file.SetCellValue(sheetName, fmt.Sprintf("K%d", rowIndex), colValueCustomer)
					case "nama":
						file.SetCellValue(sheetName, fmt.Sprintf("L%d", rowIndex), colValueCustomer)
					default:
						log.Println("Key tidak diketahui")
					}

				}
				rowIndex++

			}
		}
	}
	// set autofilter
	if err := file.AutoFilter(sheetName, "A2:K2", []excelize.AutoFilterOptions{}); err != nil {
		log.Fatal("Error", err.Error())
	}

	file.SetActiveSheet(0)

	// Save excel file
	localUploadDir := "./uploads"
	localFilePath := filepath.Join(localUploadDir, fileName)

	if err := file.SaveAs(localFilePath); err != nil {
		log.Println(err.Error())
	}

	// upload to s3
	if err := aws.NewConnect().S3.UploadFile(bucketExport, localFilePath, filePath); err != nil {
		log.Println("failed upload to s3:", err.Error())
		return err
	}

	if err := os.Remove(localFilePath); err != nil {
		log.Println("failed to remove uploaded file:", err.Error())
		return err
	} else {
		log.Println("file removed successfully:", localFilePath)
	}
	return nil
}
