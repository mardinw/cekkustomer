package dtos

type TTLSessionData struct {
	UUID        string `dynamodbav:"uuid"`
	AccessToken string `dynamodbav:"access_token"`
	CreatedAt   int64  `dynamodbav:"created_at"`
	ExpireAt    int64  `dynamodbav:"expire_at"`
}
