package dtos

type ImportCustomerXls struct {
	Agencies       string `dynamodb:"agencies"`
	CardNumber     int64  `dynamodb:"card_number"`
	FirstName      string `dynamodb:"first_name"`
	Address3       string `dynamodb:"address_3"`
	Address4       string `dynamodb:"address_4"`
	ZipCode        string `dynamodb:"home_zip_code"`
	ConcatCustomer string `dynamodb:"concat_customer"`
	Collector      string `dynamodb:"collector"`
	Created        int64  `dynamodb:"created"`
}
