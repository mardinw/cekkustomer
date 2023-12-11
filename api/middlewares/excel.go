package middlewares

import (
	"errors"
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
