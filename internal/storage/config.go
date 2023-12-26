package storagecfg

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log/slog"
	"os"
)

type Config struct {
	StorageType string `yaml:"storage_type"`
	ConnType    string `yaml:"conn_type"`
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	DB          string `yaml:"db"`
}

func New(log *slog.Logger) *Config {
	const pathOfConfig = `./storage/config/config.yaml`
	if _, err := os.Stat(pathOfConfig); os.IsNotExist(err) {
		log.Error(" There is no config file", slog.String("path", pathOfConfig), slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})
	}
	var cfg Config
	if err := cleanenv.ReadConfig(pathOfConfig, &cfg); err != nil {
		log.Error("Cannot read config", slog.String("path", pathOfConfig), slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})
	}
	return &cfg
}
