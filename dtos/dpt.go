package dtos

type DataPreview struct {
	CardNumber int64  `json:"card_number"`
	FirstName  string `json:"first_name"`
	Collector  string `json:"collector"`
	Agencies   string `json:"agencies,omitempty"`
	Address3   string `json:"address_3"`
	Address4   string `json:"address_4"`
	ZipCode    string `json:"home_zip_code"`
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
