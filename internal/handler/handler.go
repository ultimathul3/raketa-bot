package handler

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vanyaio/raketa-bot/internal/service"
	"github.com/vanyaio/raketa-bot/internal/storage"
)

type Handler struct {
	srv  service.Service
	bot  *tgbotapi.BotAPI
	strg storage.Storage
}

func NewHandler(srv service.Service, bot *tgbotapi.BotAPI, strg storage.Storage) *Handler {
	return &Handler{
		srv:  srv,
		bot:  bot,
		strg: strg,
	}
}

func (h *Handler) HandleUpdates(config tgbotapi.UpdateConfig) {
	updates := h.bot.GetUpdatesChan(config)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		userID := update.Message.From.ID
		chatID := update.Message.Chat.ID
		text := update.Message.Text

		// URL input
		if h.strg.GetState(userID) == storage.UrlInput {
			if _, err := url.ParseRequestURI(text); err != nil {
				msg := tgbotapi.NewMessage(chatID, "Invalid URL")
				h.bot.Send(msg)
				continue
			}
			h.strg.GetCallback(userID, storage.UrlInput)(text)
			h.strg.SetState(userID, storage.Menu, nil)
			continue
		}

		// /start handle
		if text == "/start" {
			err := h.srv.SignUp(context.Background(), userID)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			msg := tgbotapi.NewMessage(
				chatID,
				fmt.Sprintf("User with id '%d' signed up", userID),
			)
			msg.ReplyMarkup = menuKeyboard
			h.bot.Send(msg)
			// Create task handle
		} else if strings.HasPrefix(text, "Create task") {
			h.strg.SetState(userID, storage.UrlInput, func(ctx ...any) {
				url := ctx[0].(string)
				err := h.srv.CreateTask(context.Background(), url)
				if err != nil {
					log.Println(err.Error())
					msg := tgbotapi.NewMessage(chatID, "Task already exists")
					h.bot.Send(msg)
					return
				}
				msg := tgbotapi.NewMessage(chatID, "Task was created")
				h.bot.Send(msg)
			})
			// Delete task handle
		} else if strings.Contains(text, "Delete task") {
			h.strg.SetState(userID, storage.UrlInput, func(ctx ...any) {
				url := ctx[0].(string)
				err := h.srv.DeleteTask(context.Background(), url)
				if err != nil {
					log.Println(err.Error())
					msg := tgbotapi.NewMessage(chatID, "Task not found")
					h.bot.Send(msg)
					return
				}
				msg := tgbotapi.NewMessage(chatID, "Task was deleted")
				h.bot.Send(msg)
			})
			// Assign worker handle
			// TODO
		} else if strings.Contains(text, "Assign worker") {
			err := h.srv.AssignUser(context.Background(), "url", userID)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			msg := tgbotapi.NewMessage(chatID, "Worker was assigned")
			h.bot.Send(msg)
			// Close task handle
		} else if strings.Contains(text, "Close task") {
			h.strg.SetState(userID, storage.UrlInput, func(ctx ...any) {
				url := ctx[0].(string)
				err := h.srv.CloseTask(context.Background(), url)
				if err != nil {
					log.Println(err.Error())
					return
				}
				msg := tgbotapi.NewMessage(chatID, "Task was closed")
				h.bot.Send(msg)
			})
			// Get open tasks handle
		} else if strings.Contains(text, "Get open tasks") {
			tasks, err := h.srv.GetOpenTasks(context.Background())
			if err != nil {
				log.Println(err.Error())
				continue
			}
			var msg tgbotapi.MessageConfig
			if tasks == nil {
				msg = tgbotapi.NewMessage(chatID, "Empty tasks list")
			} else {
				msg = tgbotapi.NewMessage(chatID, "Current open tasks:")
				msg.ReplyMarkup = NewTasksKeyboard(tasks)
			}
			h.bot.Send(msg)
		}

		if h.strg.GetState(userID) == storage.UrlInput {
			msg := tgbotapi.NewMessage(chatID, "Enter task URL")
			h.bot.Send(msg)
		}
	}
}
