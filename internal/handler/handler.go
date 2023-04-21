package handler

import (
	"context"
	"fmt"
	"strings"

	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vanyaio/raketa-bot/internal/service"
)

type Handler struct {
	srv service.Service
	bot *tgbotapi.BotAPI
}

func NewHandler(srv service.Service, bot *tgbotapi.BotAPI) *Handler {
	return &Handler{
		srv: srv,
		bot: bot,
	}
}

func (h *Handler) HandleUpdates(config tgbotapi.UpdateConfig) {
	updates := h.bot.GetUpdatesChan(config)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		userID := update.Message.From.ID

		if update.Message.Text == "/start" {
			err := h.srv.SignUp(context.Background(), userID)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				fmt.Sprintf("User with id '%d' signed up", userID),
			)
			msg.ReplyMarkup = menuKeyboard
			h.bot.Send(msg)
			continue
		}

		if strings.HasPrefix(update.Message.Text, "Create task") {
			err := h.srv.CreateTask(context.Background(), "url")
			if err != nil {
				log.Println(err.Error())
				continue
			}
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"Task was created",
			)
			h.bot.Send(msg)
			continue
		}

		if strings.Contains(update.Message.Text, "Delete task") {
			err := h.srv.DeleteTask(context.Background(), "url")
			if err != nil {
				log.Println(err.Error())
				continue
			}
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"Task was deleted",
			)
			h.bot.Send(msg)
			continue
		}

		if strings.Contains(update.Message.Text, "Assign worker") {
			err := h.srv.AssignUser(context.Background(), "url", userID)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"Worker was assigned",
			)
			h.bot.Send(msg)
			continue
		}

		if strings.Contains(update.Message.Text, "Close task") {
			err := h.srv.CloseTask(context.Background(), "url")
			if err != nil {
				log.Println(err.Error())
				continue
			}
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"Task was closed",
			)
			h.bot.Send(msg)
			continue
		}

		if strings.Contains(update.Message.Text, "Get open tasks") {
			tasks, err := h.srv.GetOpenTasks(context.Background())
			if err != nil {
				log.Println(err.Error())
				continue
			}
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"Current open tasks:",
			)
			msg.ReplyMarkup = NewTasksKeyboard(tasks)
			h.bot.Send(msg)
			continue
		}
	}
}
