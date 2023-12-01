package models

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type DPT struct {
	CardNumber   int64          `json:"card_number"`
	FirstName    string         `json:"nama"`
	HomeAddress3 sql.NullString `json:"address_3"`
	HomeAddress4 sql.NullString `json:"address_4"`
	HomeZipCode  int32          `json:"zip_code"`
	Kodepos      int32          `json:"kodepos"`
	Kelurahan    string         `json:"kelurahan"`
	Kecamatan    string         `json:"kecamatan"`
}

func (dpt *DPT) GetAll(db *sql.DB) ([]DPT, error) {
	query := `
	select t1.card_number AS card_number,t1.first_name AS first_name,t1.home_address_3 AS home_address_3,t1.home_address_4 AS home_address_4,t1.home_zip_code AS home_zip_code,t2.kodepos AS kodepos,t2.kel AS kel,t2.kec AS kec from (datadpt.customer t1 join datadpt.dpt_cibeunyingkidul t2 on((t1.concat_customer = t2.concat)))
	`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	defer rows.Close()

	var result []DPT
	for rows.Next() {
		var each = DPT{}
		var err = rows.Scan(
			&each.CardNumber,
			&each.FirstName,
			&each.HomeAddress3,
			&each.HomeAddress4,
			&each.HomeZipCode,
			&each.Kodepos,
			&each.Kecamatan,
			&each.Kelurahan,
		)
		if err != nil {
			log.Println("record not found")
			return nil, err
		}

		result = append(result, each)

	}

	return result, nil
}
