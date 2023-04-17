package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Bot struct {
		Token string `env:"BOT_TOKEN"`
		Debug bool   `env:"BOT_DEBUG"`
	}

	GRPC struct {
		ServerHost string `env:"GRPC_SERVER_HOST"`
		ServerPort uint8  `env:"GRPC_SERVER_PORT"`
	}
}

func ReadEnvFile() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(".env", &cfg)

	return &cfg, err
}
