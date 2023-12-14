package middlewares

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"

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
		"concat_customer",
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

func CreateExcel(jsonData string) error {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return err
	}

	// membuat file excel

	file := excelize.NewFile()

	sheetName := "Match Data"
	headers := []string{"Card Number", "First Name", "Collector", "Agencies", "Address 3", "Address 4", "Zip Code", "Kode Pos", "Nama", "Kelurahan", "Kecamatan", "Lokasi"}
	file.SetSheetName(file.GetSheetName(0), sheetName)
	file.SetCellValue(sheetName, "A1", "Data Customer")
	file.SetCellValue(sheetName, "H1", "Data Match")
	file.MergeCell(sheetName, "A1", "G1")
	file.MergeCell(sheetName, "H1", "L1")

	// Membuat table header untuk data
	for colIndex, header := range headers {
		colName, _ := excelize.ColumnNumberToName(colIndex + 1)
		file.SetCellValue(sheetName, fmt.Sprintf("%s2", colName), header)
	}

	if err := file.AutoFilter(sheetName, "A2:L2", []excelize.AutoFilterOptions{}); err != nil {
		log.Fatal("Error", err.Error())
	}

	// tambahkan data ke sheet
	rowIndex := 3
	for tableName, tableData := range data {
		for _, rowData := range tableData.(map[string]interface{})["db_match"].([]interface{}) {
			for colIndex, colValue := range rowData.(map[string]interface{}) {
				colIndexConv, err := strconv.Atoi(colIndex)
				if err != nil {
					log.Println(err.Error())
				}
				colName, _ := excelize.ColumnNumberToName(colIndexConv + 1)
				file.SetCellValue(sheetName, fmt.Sprintf("%s%d", colName, rowIndex), colValue)
			}

			// tambah table nama
			file.SetCellValue(sheetName, fmt.Sprintf("A%d", rowIndex), tableName)
			rowIndex++
		}
	}

	// set autofilter
	if err := file.AutoFilter(sheetName, "A2:L2", []excelize.AutoFilterOptions{}); err != nil {
		log.Fatal("Error", err.Error())
	}

	file.SetActiveSheet(0)

	if err := file.SaveAs("./file1.xlsx"); err != nil {
		log.Println(err.Error())
	}

	return nil
}
