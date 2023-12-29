package postgres

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log/slog"
	"os"
)

type Config struct {
	ConnType   string `yaml:"conn_type"`
	ConnString string `yaml:"conn_string"`
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	DB         string `yaml:"db"`
}

func NewConfig(log *slog.Logger) *Config {
	const pathOfConfig = "./internal/db/config.yaml"
	if _, err := os.Stat(pathOfConfig); os.IsNotExist(err) {
		log.Error("there is no config file ", slog.String("file", pathOfConfig))
		panic("cant create config")
	}
	var cfg Config
	if err := cleanenv.ReadConfig(pathOfConfig, &cfg); err != nil {
		log.Error("cannot read config", slog.String("error", err.Error()))
		panic("cant create config")
	}
	return &cfg
}
