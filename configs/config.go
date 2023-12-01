package configs

type AppConfiguration struct {
	Mode     string `env:"GIN_MODE"`
	Port     int    `env:"PORT"`
	AppEnv   string `env:"APP_ENV"`
	Version  string `env:"VERSION"`
	Database DBConfig
}

type DBConfig struct {
	Type         string `env:"DB_TYPE,default=postgres"`
	EndPoint     string `env:"DB_ENDPOINT"`
	ReadEndPoint string `env:"READ_ENDPOINT"`
	Name         string `env:"DB_NAME,default=postgress"`
	User         string `env:"DB_USER,default=postgress"`
	Password     string `env:"DB_PASSWORD"`
}
