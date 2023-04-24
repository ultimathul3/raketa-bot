package handler

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	raketapb "github.com/vanyaio/raketa-backend/proto"
	"github.com/vanyaio/raketa-bot/internal/types"
)

type storage interface {
	GetState(ID int64) types.State
	GetCallback(ID int64) types.Callback
	SetState(ID int64, state types.State, callback types.Callback)
}

type service interface {
	SignUp(ctx context.Context, id int64) error
	CreateTask(ctx context.Context, url string) error
	DeleteTask(ctx context.Context, url string) error
	AssignUser(ctx context.Context, url string, userID int64) error
	CloseTask(ctx context.Context, url string) error
	GetOpenTasks(ctx context.Context) ([]*raketapb.Task, error)
}

type Handler struct {
	srv     service
	bot     *tgbotapi.BotAPI
	storage storage
}

func NewHandler(srv service, bot *tgbotapi.BotAPI, storage storage) *Handler {
	return &Handler{
		srv:     srv,
		bot:     bot,
		storage: storage,
	}
}

func (h *Handler) HandleUpdates(ctx context.Context, config tgbotapi.UpdateConfig) {
	updates := h.bot.GetUpdatesChan(config)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		userID := update.Message.From.ID
		chatID := update.Message.Chat.ID
		text := update.Message.Text

		// URL input
		if h.storage.GetState(userID) == types.UrlInput {
			if _, err := url.ParseRequestURI(text); err != nil {
				msg := tgbotapi.NewMessage(chatID, "Invalid URL")
				h.bot.Send(msg)
				continue
			}
			h.storage.GetCallback(userID)(text)
			h.storage.SetState(userID, types.Menu, nil)
			continue
		}

		// /start handle
		if text == "/start" {
			err := h.srv.SignUp(ctx, userID)
			if err != nil {
				log.Println(err.Error())
			}
			msg := tgbotapi.NewMessage(
				chatID,
				fmt.Sprintf("User with id '%d' signed up", userID),
			)
			msg.ReplyMarkup = menuKeyboard
			h.bot.Send(msg)
			// Create task handle
		} else if strings.HasPrefix(text, "Create task") {
			h.storage.SetState(userID, types.UrlInput, func(url string) {
				err := h.srv.CreateTask(ctx, url)
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
			h.storage.SetState(userID, types.UrlInput, func(url string) {
				err := h.srv.DeleteTask(ctx, url)
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
			err := h.srv.AssignUser(ctx, "url", userID)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			msg := tgbotapi.NewMessage(chatID, "Worker was assigned")
			h.bot.Send(msg)
			// Close task handle
		} else if strings.Contains(text, "Close task") {
			h.storage.SetState(userID, types.UrlInput, func(url string) {
				err := h.srv.CloseTask(ctx, url)
				if err != nil {
					log.Println(err.Error())
					return
				}
				msg := tgbotapi.NewMessage(chatID, "Task was closed")
				h.bot.Send(msg)
			})
			// Get open tasks handle
		} else if strings.Contains(text, "Get open tasks") {
			tasks, err := h.srv.GetOpenTasks(ctx)
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

		if h.storage.GetState(userID) == types.UrlInput {
			msg := tgbotapi.NewMessage(chatID, "Enter task URL")
			h.bot.Send(msg)
		}
	}
}
