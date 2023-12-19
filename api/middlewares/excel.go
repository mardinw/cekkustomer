package middlewares

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/xuri/excelize/v2"
)

type MapCustomer []map[string]interface{}

func containsKey(keys []string, key string) bool {
	for _, k := range keys {
		if k == key {
			return true
		}
	}

	return false
}

func ReadExcel(fileName io.ReadCloser) (MapCustomer, error) {
	xlsx, err := excelize.OpenReader(fileName)
	if err != nil {
		log.Println(err.Error())
	}

	sheetName := xlsx.GetSheetList()[0]

	rows, err := xlsx.GetRows(sheetName)
	if err != nil {
		log.Println(err.Error())
	}

	// konversi data excel to json
	var result MapCustomer
	keys := rows[0]

	checkKeys := []string{
		"card_number",
		"first_name",
		"address_3",
		"address_4",
		"home_zip_code",
		"collector",
		"concat_customer (nama + tgl lahir)",
	}

	// mengecek table header
	for _, key := range checkKeys {
		if containsKey(keys, key) {
			log.Println("pass")
		} else {
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

func CreateExcel(jsonData, bucketExport, filePath string) error {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return err
	}

	// membuat file excel

	file := excelize.NewFile()

	sheetName := "Match Data"
	headers := []string{"Card Number", "First Name", "Collector", "Address 3", "Address 4", "Zip Code", "Kode Pos", "Nama", "Kelurahan", "Kecamatan", "Lokasi"}
	file.SetSheetName(file.GetSheetName(0), sheetName)
	file.SetCellValue(sheetName, "A1", "Data Customer")
	file.SetCellValue(sheetName, "G1", "Data Match")
	file.MergeCell(sheetName, "A1", "F1")
	file.MergeCell(sheetName, "G1", "K1")

	// Membuat table header untuk data
	for colIndex, header := range headers {
		colName, _ := excelize.ColumnNumberToName(colIndex + 1)
		file.SetCellValue(sheetName, fmt.Sprintf("%s2", colName), header)
	}

	// tambahkan data ke sheet
	rowIndex := 3
	for tableName, tableData := range data {
		// data customer
		for colName, colValueCustomer := range tableData.(map[string]interface{}) {
			switch colName {
			case "card_number":
				file.SetCellValue(sheetName, fmt.Sprintf("A%d", rowIndex), colValueCustomer)
			case "first_name":
				file.SetCellValue(sheetName, fmt.Sprintf("B%d", rowIndex), colValueCustomer)
			case "collector":
				file.SetCellValue(sheetName, fmt.Sprintf("C%d", rowIndex), colValueCustomer)
			case "address_3":
				file.SetCellValue(sheetName, fmt.Sprintf("D%d", rowIndex), colValueCustomer)
			case "address_4":
				file.SetCellValue(sheetName, fmt.Sprintf("E%d", rowIndex), colValueCustomer)
			case "home_zip_code":
				file.SetCellValue(sheetName, fmt.Sprintf("F%d", rowIndex), colValueCustomer)
			default:
				log.Println("Key tidak diketahui")
			}

		}

		// ambil data match
		for _, rowData := range tableData.(map[string]interface{})["db_match"].([]interface{}) {
			for colName, colValue := range rowData.(map[string]interface{}) {
				switch colName {
				case "kodepos":
					file.SetCellValue(sheetName, fmt.Sprintf("G%d", rowIndex), colValue)
				case "nama":
					file.SetCellValue(sheetName, fmt.Sprintf("H%d", rowIndex), colValue)
				case "kelurahan":
					file.SetCellValue(sheetName, fmt.Sprintf("I%d", rowIndex), colValue)
				case "kecamatan":
					file.SetCellValue(sheetName, fmt.Sprintf("J%d", rowIndex), colValue)
				default:
					log.Println("Key tidak dikenali")
				}
			}

			// tambah lokasi
			file.SetCellValue(sheetName, fmt.Sprintf("K%d", rowIndex), tableName)
			rowIndex++
		}
	}

	// set autofilter
	if err := file.AutoFilter(sheetName, "A2:K2", []excelize.AutoFilterOptions{}); err != nil {
		log.Fatal("Error", err.Error())
	}

	file.SetActiveSheet(0)

	if err := file.SaveAs("./file1.xlsx"); err != nil {
		log.Println(err.Error())
	}

	return nil
}
