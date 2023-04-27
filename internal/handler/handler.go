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
	GetState(userID int64) types.State
	GetData(userID int64, key string) any
	SetState(userID int64, state types.State)
	SetStateWithData(userID int64, state types.State, key string, data any)
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
		state := h.storage.GetState(userID)

		var msg tgbotapi.MessageConfig

		switch text {
		case startCommand:
			if err := h.srv.SignUp(ctx, userID); err != nil {
				handleError(err, h.bot, chatID)
				continue
			}
			msg = tgbotapi.NewMessage(chatID, getUserSignedUpMessage(userID))
			msg.ReplyMarkup = menuKeyboard

		case createTaskCommand:
			msg = tgbotapi.NewMessage(chatID, enterTaskUrlMessage)
			h.storage.SetState(userID, types.CreateTaskUrlInput)

		case deleteTaskCommand:
			msg = tgbotapi.NewMessage(chatID, enterTaskUrlMessage)
			h.storage.SetState(userID, types.DeleteTaskUrlInput)

		case assignWorkerCommand:
			msg = tgbotapi.NewMessage(chatID, enterTaskUrlMessage)
			h.storage.SetState(userID, types.AssignWorkerUrlInput)

		case closeTaskCommand:
			msg = tgbotapi.NewMessage(chatID, enterTaskUrlMessage)
			h.storage.SetState(userID, types.CloseTaskUrlInput)

		case getOpenTasksCommand:
			tasks, err := h.srv.GetOpenTasks(ctx)
			if err != nil {
				handleError(err, h.bot, chatID)
				continue
			}
			if tasks == nil {
				msg = tgbotapi.NewMessage(chatID, emptyTasksListMessage)
			} else {
				msg = tgbotapi.NewMessage(chatID, currentOpenTasksMessage)
				msg.ReplyMarkup = NewTasksKeyboard(tasks)
			}

		default:

			switch state {
			case types.CreateTaskUrlInput:
				url, err := handleUrlInput(h.bot, chatID, text)
				if err != nil {
					continue
				}
				if err := h.srv.CreateTask(ctx, url); err != nil {
					handleError(err, h.bot, chatID)
					continue
				}
				msg = tgbotapi.NewMessage(chatID, taskWasCreatedMessage)
				h.storage.SetState(userID, types.Menu)

			case types.DeleteTaskUrlInput:
				url, err := handleUrlInput(h.bot, chatID, text)
				if err != nil {
					continue
				}
				if err := h.srv.DeleteTask(ctx, url); err != nil {
					handleError(err, h.bot, chatID)
					continue
				}
				msg = tgbotapi.NewMessage(chatID, taskWasDeletedMessage)
				h.storage.SetState(userID, types.Menu)

			case types.AssignWorkerUrlInput:
				url, err := handleUrlInput(h.bot, chatID, text)
				if err != nil {
					continue
				}
				msg = tgbotapi.NewMessage(chatID, enterUserIdMessage)
				h.storage.SetStateWithData(userID, types.AssignWorkerIdInput, "url", url)

			case types.AssignWorkerIdInput:
				id, err := handleIdInput(h.bot, chatID, text)
				if err != nil {
					continue
				}
				url := h.storage.GetData(userID, "url").(string)
				if err := h.srv.AssignUser(ctx, url, id); err != nil {
					handleError(err, h.bot, chatID)
					continue
				}
				msg = tgbotapi.NewMessage(chatID, workerWasAssignedMessage)
				h.storage.SetState(userID, types.Menu)

			case types.CloseTaskUrlInput:
				url, err := handleUrlInput(h.bot, chatID, text)
				if err != nil {
					continue
				}
				if err := h.srv.CloseTask(ctx, url); err != nil {
					handleError(err, h.bot, chatID)
					continue
				}
				msg = tgbotapi.NewMessage(chatID, taskWasClosedMessage)
				h.storage.SetState(userID, types.Menu)
			}
		}

		h.bot.Send(msg)
	}
}

func handleError(err error, bot *tgbotapi.BotAPI, chatID int64) {
	log.Println(err.Error())
	msg := tgbotapi.NewMessage(chatID, err.Error())
	bot.Send(msg)
}

func handleUrlInput(bot *tgbotapi.BotAPI, chatID int64, input string) (string, error) {
	_, err := url.ParseRequestURI(input)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, invalidUrlMessage)
		bot.Send(msg)
	}

	return input, err
}

func handleIdInput(bot *tgbotapi.BotAPI, chatID int64, input string) (int64, error) {
	id, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, invalidUserIdMessage)
		bot.Send(msg)
	}

	return id, err
}

func getUserSignedUpMessage(userID int64) string {
	return fmt.Sprintf(userSignedUpMessageFmt, userID)
}
