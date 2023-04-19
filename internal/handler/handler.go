package handler

import (
	"context"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vanyaio/raketa-bot/internal/service"
)

type Handler struct {
	srv service.Service
}

func NewHandler(srv service.Service) *Handler {
	return &Handler{
		srv: srv,
	}
}

func (h *Handler) HandleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message != nil {
			userID := update.Message.From.ID
			textParts := strings.Split(update.Message.Text, " ")
			command := textParts[0]

			switch command {
			case "/start":
				h.srv.SignUp(context.Background(), userID)
			case "/createtask":
				url := textParts[1]
				h.srv.CreateTask(context.Background(), url)
			case "/deletetask":
				url := textParts[1]
				h.srv.DeleteTask(context.Background(), url)
			case "/assignworker":
				url := textParts[1]
				h.srv.AssignWorker(context.Background(), url, userID)
			case "/closetask":
				url := textParts[1]
				h.srv.CloseTask(context.Background(), url)
			case "/getopentasks":
				h.srv.GetOpenTasks(context.Background())
			}
		}
	}
}
