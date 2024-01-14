package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string `yaml:"env"`
	Admin      `yaml:"admin"`
	HTTPServer `yaml:"http_server"`
	DB         `yaml:"db"`
}

type Admin struct {
	Login    string `yaml:"login"`
	Password string
}

type DB struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"auth"`
	Password string
}

type HTTPServer struct {
	Address     string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

func New() *Config {
	const pathOfConfig = `./config/config.yaml`
	if _, err := os.Stat(pathOfConfig); os.IsNotExist(err) {
		log.Fatalf("there is no config file in %s", pathOfConfig)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(pathOfConfig, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("cannot load passwords %s", err)
	}
	cfg.Admin.Password = os.Getenv("ADMIN_PASSWORD")
	cfg.DB.Password = os.Getenv("DB_PASSWORD")
	return &cfg
}
