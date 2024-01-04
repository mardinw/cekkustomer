package models

type ImportCustomerXls struct {
	CardNumber     string `json:"card_number"`
	NIK            string `json:"nik"`
	NIKCheck       string `json:"nik_check"`
	FirstName      string `json:"first_name"`
	Collector      string `json:"collector"`
	Agencies       string `json:"agencies"`
	Address3       string `json:"address_3"`
	Address4       string `json:"address_4"`
	ZipCode        string `json:"home_zip_code"`
	ConcatCustomer string `json:"concat_customer"`
	Files          string `json:"files"`
	Created        int64  `json:"created"`
}
