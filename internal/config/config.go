package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	BotToken string `env:"BOT_TOKEN"`
	BotDebug bool   `env:"BOT_DEBUG"`
}

func ReadEnvFile() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(".env", &cfg)

	return &cfg, err
}
