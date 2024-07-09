package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
)

type Config struct {
	TokenTTL time.Duration `yaml:"token_ttl" env-required:"true"`
	App      AppConfig
	Postgres PostgresConfig
}

type AppConfig struct {
	Port uint16 `yaml:"port"`
	Addr string `yaml:"addr"`
}

type PostgresConfig struct {
	Addr     string `yaml:"addr"`
	Port     uint16 `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DB       string `yaml:"db"`
}

// TODO: comment fuction
func MustLoad() *Config {
	var cfg Config

	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic(errors.Wrap(err, "config file dost not exists"))
	}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic(errors.Wrap(err, "failed to read config"))
	}
	return &cfg
}

// fetchConfigPath fetches config path from command line flar or env variable
// Priority: flag > env > default
// Default value is empty string
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path fo config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
