package configs

type AppConfiguration struct {
	Mode    string `env:"GIN_MODE"`
	Port    int    `env:"PORT"`
	AppEnv  string `env:"APP_ENV"`
	Version string `env:"VERSION"`
}
