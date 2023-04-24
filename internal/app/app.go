package app

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vanyaio/raketa-bot/internal/config"
	"github.com/vanyaio/raketa-bot/internal/handler"
	"github.com/vanyaio/raketa-bot/internal/service"
	"github.com/vanyaio/raketa-bot/internal/storage"
	"github.com/vanyaio/raketa-bot/pkg/client"
)

func Run(cfg *config.Config) {
	bot, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = cfg.Bot.Debug

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	client, err := client.NewGrpcClient(cfg.GRPC.ServerHost, cfg.GRPC.ServerPort)
	if err != nil {
		log.Fatal(err)
	}

	storage := storage.NewStateStorage()
	service := service.NewRaketaService(client)
	handler := handler.NewHandler(service, bot, storage)
	handler.HandleUpdates(context.Background(), u)
}
