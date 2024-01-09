package database

type Config struct {
	Driver Driver `mapstructure:"driver" validate:"required" default:"cockroachdb"`
	URI    string `mapstructure:"uri" validate:"required" default:"postgres://root@localhost:26257/defaultdb"`
}
