package dtos

type Customer struct {
	CardNumber     int64  `json:"card_number"`
	FirstName      string `json:"first_name"`
	Address3       string `json:"address_3"`
	Address4       string `json:"address_4"`
	ZipCode        string `json:"home_zip_code"`
	ConcatCustomer string `json:"concat_customer"`
	Collector      string `json:"collector"`
	Agencies       string `json:"agencies"`
}
