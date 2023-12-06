package middlewares

import (
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
)

type MapCustomer map[string]interface{}

func ReadExcel(fileName string) {
	file, err := excelize.OpenFile(fileName)
	if err != nil {
		log.Println(err.Error())
	}

	defer file.Close()

	sheets := file.GetSheetList()

	for _, sheet := range sheets {
		rows, err := file.GetRows(sheet)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		for rowIndex, row := range rows {
			for colIndex, cell := range row {
				log.Printf("Sheet: %s, Row: %d, Col: %d, Value: %s\n", sheet, rowIndex+1, colIndex+1, cell)
			}
		}
	}
}
