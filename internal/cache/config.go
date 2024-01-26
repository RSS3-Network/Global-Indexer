package cache

type Config struct {
	Endpoints string `mapstructure:"endpoints" validate:"required" default:"127.0.0.1:6379"`
	Password  string `mapstructure:"password"`
	Username  string `mapstructure:"username"`
	DB        int    `mapstructure:"db"`
}
