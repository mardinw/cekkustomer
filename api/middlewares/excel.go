package middlewares

import (
	"io"
	"log"

	"github.com/xuri/excelize/v2"
)

type MapCustomer []map[string]interface{}

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
