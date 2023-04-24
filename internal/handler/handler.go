package handler

import (
	"context"
	"fmt"
	"net/url"

	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	GetOpenTasks(ctx context.Context) ([]types.Task, error)
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

		if h.storage.GetState(userID) == types.UrlInput {
			if _, err := url.ParseRequestURI(text); err != nil {
				msg := tgbotapi.NewMessage(chatID, invalidUrlMessage)
				h.bot.Send(msg)
				continue
			}
			h.storage.GetCallback(userID)(text)
			h.storage.SetState(userID, types.Menu, nil)
			continue
		}

		switch text {
		case startCommand:
			err := h.srv.SignUp(ctx, userID)
			var msg tgbotapi.MessageConfig
			if err != nil {
				log.Println(err.Error())
				msg = tgbotapi.NewMessage(chatID, err.Error())
			} else {
				msg = tgbotapi.NewMessage(
					chatID,
					fmt.Sprintf(userSignedUpMessage, userID),
				)
				msg.ReplyMarkup = menuKeyboard
			}
			h.bot.Send(msg)
		case createTaskCommand:
			h.storage.SetState(userID, types.UrlInput, func(url string) {
				err := h.srv.CreateTask(ctx, url)
				if err != nil {
					log.Println(err.Error())
					msg := tgbotapi.NewMessage(chatID, taskAlreadyExistsMessage)
					h.bot.Send(msg)
					return
				}
				msg := tgbotapi.NewMessage(chatID, taskWasCreatedMessage)
				h.bot.Send(msg)
			})
		case deleteTaskCommand:
			h.storage.SetState(userID, types.UrlInput, func(url string) {
				err := h.srv.DeleteTask(ctx, url)
				if err != nil {
					log.Println(err.Error())
					msg := tgbotapi.NewMessage(chatID, taskNotFoundMessage)
					h.bot.Send(msg)
					return
				}
				msg := tgbotapi.NewMessage(chatID, taskWasDeletedMessage)
				h.bot.Send(msg)
			})
		case assignWorkerCommand:
			// TODO
			err := h.srv.AssignUser(ctx, "url", userID)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			msg := tgbotapi.NewMessage(chatID, workerWasAssignedMessage)
			h.bot.Send(msg)
		case closeTaskCommand:
			h.storage.SetState(userID, types.UrlInput, func(url string) {
				err := h.srv.CloseTask(ctx, url)
				if err != nil {
					log.Println(err.Error())
					return
				}
				msg := tgbotapi.NewMessage(chatID, taskWasClosedMessage)
				h.bot.Send(msg)
			})
		case getOpenTasksCommand:
			tasks, err := h.srv.GetOpenTasks(ctx)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			var msg tgbotapi.MessageConfig
			if tasks == nil {
				msg = tgbotapi.NewMessage(chatID, emptyTasksListMessage)
			} else {
				msg = tgbotapi.NewMessage(chatID, currentOpenTasksMessage)
				msg.ReplyMarkup = NewTasksKeyboard(tasks)
			}
			h.bot.Send(msg)
		}

		if h.storage.GetState(userID) == types.UrlInput {
			msg := tgbotapi.NewMessage(chatID, enterTaskUrlMessage)
			h.bot.Send(msg)
		}
	}
}
