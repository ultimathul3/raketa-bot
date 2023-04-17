package main

import (
	"log"

	"github.com/vanyaio/raketa-bot/internal/app"
	"github.com/vanyaio/raketa-bot/internal/config"
)

func main() {
	cfg, err := config.ReadEnvFile()
	if err != nil {
		log.Fatal(err)
	}

	app.Run(cfg)
}
