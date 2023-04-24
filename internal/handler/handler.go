package handler

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

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

		switch h.storage.GetState(userID) {
		case types.UrlInput:
			if _, err := url.ParseRequestURI(text); err != nil {
				msg := tgbotapi.NewMessage(chatID, invalidUrlMessage)
				h.bot.Send(msg)
				continue
			}
			h.storage.GetCallback(userID)(text)
			continue

		case types.IdInput:
			if _, err := strconv.ParseInt(text, 10, 64); err != nil {
				msg := tgbotapi.NewMessage(chatID, invalidUserIdMessage)
				h.bot.Send(msg)
				continue
			}
			h.storage.GetCallback(userID)(text)
			continue
		}

		switch text {
		case startCommand:
			if err := h.srv.SignUp(ctx, userID); err != nil {
				handleError(err, h.bot, chatID)
				continue
			}
			msg := tgbotapi.NewMessage(
				chatID,
				fmt.Sprintf(userSignedUpMessage, userID),
			)
			msg.ReplyMarkup = menuKeyboard
			h.bot.Send(msg)

		case createTaskCommand:
			msg := tgbotapi.NewMessage(chatID, enterTaskUrlMessage)
			h.bot.Send(msg)
			h.storage.SetState(userID, types.UrlInput, func(url string) {
				if err := h.srv.CreateTask(ctx, url); err != nil {
					handleError(err, h.bot, chatID)
					return
				}
				msg := tgbotapi.NewMessage(chatID, taskWasCreatedMessage)
				h.bot.Send(msg)
				h.storage.SetState(userID, types.Menu, nil)
			})

		case deleteTaskCommand:
			msg := tgbotapi.NewMessage(chatID, enterTaskUrlMessage)
			h.bot.Send(msg)
			h.storage.SetState(userID, types.UrlInput, func(url string) {
				if err := h.srv.DeleteTask(ctx, url); err != nil {
					handleError(err, h.bot, chatID)
					return
				}
				msg := tgbotapi.NewMessage(chatID, taskWasDeletedMessage)
				h.bot.Send(msg)
				h.storage.SetState(userID, types.Menu, nil)
			})

		case assignWorkerCommand:
			msg := tgbotapi.NewMessage(chatID, enterTaskUrlMessage)
			h.bot.Send(msg)
			h.storage.SetState(userID, types.UrlInput, func(url string) {
				msg := tgbotapi.NewMessage(chatID, enterUserIdMessage)
				h.bot.Send(msg)
				h.storage.SetState(userID, types.IdInput, func(idInput string) {
					id, _ := strconv.ParseInt(idInput, 10, 64)
					if err := h.srv.AssignUser(ctx, url, id); err != nil {
						handleError(err, h.bot, chatID)
						return
					}
					msg := tgbotapi.NewMessage(chatID, workerWasAssignedMessage)
					h.bot.Send(msg)
					h.storage.SetState(userID, types.Menu, nil)
				})
			})

		case closeTaskCommand:
			msg := tgbotapi.NewMessage(chatID, enterTaskUrlMessage)
			h.bot.Send(msg)
			h.storage.SetState(userID, types.UrlInput, func(url string) {
				if err := h.srv.CloseTask(ctx, url); err != nil {
					handleError(err, h.bot, chatID)
					return
				}
				msg := tgbotapi.NewMessage(chatID, taskWasClosedMessage)
				h.bot.Send(msg)
				h.storage.SetState(userID, types.Menu, nil)
			})

		case getOpenTasksCommand:
			tasks, err := h.srv.GetOpenTasks(ctx)
			if err != nil {
				handleError(err, h.bot, chatID)
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
	}
}

func handleError(err error, bot *tgbotapi.BotAPI, chatID int64) {
	log.Println(err.Error())
	msg := tgbotapi.NewMessage(chatID, err.Error())
	bot.Send(msg)
}
