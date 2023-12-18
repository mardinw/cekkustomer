package dtos

import "database/sql"

type DPT struct {
	CardNumber   int64          `json:"card_number"`
	Collector    string         `json:"collector"`
	FirstName    string         `json:"first_name"`
	HomeAddress3 sql.NullString `json:"address_3"`
	HomeAddress4 sql.NullString `json:"address_4"`
	HomeZipCode  int32          `json:"zip_code"`
	Kodepos      int32          `json:"kodepos"`
	Kelurahan    string         `json:"kelurahan"`
	Kecamatan    string         `json:"kecamatan"`
}

type CheckDPT struct {
	CardNumber     int64  `json:"card_number"`
	FirstName      string `json:"first_name"`
	Collector      string `json:"collector"`
	Agencies       string `json:"agencies"`
	Address3       string `json:"address_3"`
	Address4       string `json:"address_4"`
	ZipCode        string `json:"home_zip_code"`
	ConcatCustomer string `json:"concat_customer,omitempty"`
	Files          string `json:"files,omitempty"`
	Nama           string `json:"nama"`
	Kodepos        int32  `json:"kodepos"`
	Kelurahan      string `json:"kelurahan"`
	Kecamatan      string `json:"kecamatan"`
}
