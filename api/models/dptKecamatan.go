package models

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type DPT struct {
	ID             int64  `json:"id"`
	KTP            string `json:"ktp"`
	Nama           string `json:"nama"`
	Lahir          string `json:"lahir"`
	JenisKelamin   string `json:"jenis_kelamin"`
	TPS            string `json:"tps"`
	Kecamatan      string `json:"kecamatan"`
	Kelurahan      string `json:"kelurahan"`
	Kodepos        string `json:"kodepos"`
	TempatTglLahir string `json:"ttl"`
	Concat         string `json:"concat"`
	Umur           int32  `json:"umur"`
}

func (dpt *DPT) GetAll(db *sql.DB) ([]DPT, error) {
	query := `
	SELECT * FROM dpt_kiaracondong
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
			&each.ID,
			&each.KTP,
			&each.Nama,
			&each.Lahir,
			&each.JenisKelamin,
			&each.Kelurahan,
			&each.Kodepos,
			&each.TempatTglLahir,
			&each.Concat,
			&each.Umur,
		)
		if err != nil {
			log.Println("record not found")
			return nil, err
		}

		result = append(result, each)

	}

	return result, nil
}
