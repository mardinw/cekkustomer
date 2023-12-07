package dtos

type ImportXlsx struct {
	Agencies string `dynamodb:"agencies"`
	Files    string `dynamodb:"files"`
	Uploaded int64  `dynamodb:"uploaded"`
}
