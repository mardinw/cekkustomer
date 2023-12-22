package configs

type AppConfiguration struct {
	Mode          string `env:"GIN_MODE"`
	Port          int    `env:"PORT"`
	AppEnv        string `env:"APP_ENV"`
	Version       string `env:"VERSION"`
	Database      DBConfig
	AwsConf       AwsConfiguration
	CognitoConfig AwsCognitoConfig
	S3Bucket      AwsS3Bucket
	DynamoConfig  AwsDynTblConfig
}

type DBConfig struct {
	Type         string `env:"DB_TYPE,default=postgres"`
	EndPoint     string `env:"DB_ENDPOINT"`
	ReadEndPoint string `env:"READ_ENDPOINT"`
	Name         string `env:"DB_NAME,default=postgress"`
	User         string `env:"DB_USER,default=postgress"`
	Password     string `env:"DB_PASSWORD"`
}

type AwsConfiguration struct {
	AwsProfile string `env:"AWS_PROFILE"`
	AwsRegion  string `env:"AWS_REGION"`
}

type AwsCognitoConfig struct {
	CognitoClientId     string `env:"COGNITO_CLIENT_ID"`
	CognitoClientSecret string `env:"COGNITO_CLIENT_SECRET"`
	CognitoUserPoolID   string `env:"COGNITO_USER_POOL_ID"`
}

type AwsS3Bucket struct {
	ImportS3 string `env:"IMPORT_S3"`
	ExportS3 string `env:"EXPORT_S3"`
}

type AwsDynTblConfig struct {
	TTLSes string `env:"TTLSES_DYM"`
}
