package configs

type AppConfiguration struct {
	Mode     string `env:"GIN_MODE"`
	Port     int    `env:"PORT"`
	AppEnv   string `env:"APP_ENV"`
	Version  string `env:"VERSION"`
	Database DBConfig
	AwsConf  AwsConfiguration
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
	AwsProfile   string `env:"AWS_PROFILE"`
	AwsRegion    string `env:"AWS_REGION"`
	ClientId     string `env:"CLIENT_ID"`
	ClientSecret string `env:"CLIENT_SECRET"`
	UserPoolID   string `env:"USER_POOL_ID"`
}
