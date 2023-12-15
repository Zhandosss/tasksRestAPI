package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string `yaml:"env"`
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

func New() *Config {
	const pathOfConfig = `./config/config.yaml`
	if _, err := os.Stat(pathOfConfig); os.IsNotExist(err) {
		log.Fatalf(" There is no config file in %s", pathOfConfig)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(pathOfConfig, &cfg); err != nil {
		log.Fatalf("Cannot read config: %s", err)
	}
	return &cfg
}
