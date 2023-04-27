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
	GetState(userID int64) (types.State, bool)
	GetData(userID int64, key string) (any, bool)
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

type isCommandInput bool

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
		userInput := update.Message.Text
		message := tgbotapi.MessageConfig{}
		var err error

		state, ok := h.storage.GetState(userID)
		if !ok {
			state = types.Menu
			h.storage.SetState(userID, state)
		}

		switch state {
		case types.Menu:
			message, _, err = h.handleCommandInput(ctx, userInput, userID, chatID)
			if err != nil {
				message = h.handleError(err, chatID)
			}

		case types.CreateTaskUrlInput:
			message, err = h.handleCreateTaskUrlInput(ctx, userInput, userID, chatID)
			if err != nil {
				message = h.handleError(err, chatID)
			}

		case types.DeleteTaskUrlInput:
			message, err = h.handleDeleteTaskUrlInput(ctx, userInput, userID, chatID)
			if err != nil {
				message = h.handleError(err, chatID)
			}

		case types.AssignWorkerUrlInput:
			message, err = h.handleAssignWorkerUrlInput(ctx, userInput, userID, chatID)
			if err != nil {
				message = h.handleError(err, chatID)
			}

		case types.AssignWorkerIdInput:
			message, err = h.handleAssignWorkerIdInput(ctx, userInput, userID, chatID)
			if err != nil {
				message = h.handleError(err, chatID)
			}

		case types.CloseTaskUrlInput:
			message, err = h.handleCloseTaskUrlInput(ctx, userInput, userID, chatID)
			if err != nil {
				message = h.handleError(err, chatID)
			}
		}

		if message.Text != "" {
			h.bot.Send(message)
		}
	}
}

func (h *Handler) handleCreateTaskUrlInput(ctx context.Context, input string, userID, chatID int64) (tgbotapi.MessageConfig, error) {
	msg, isCommandInput, err := h.handleCommandInput(ctx, input, userID, chatID)
	if err != nil {
		return tgbotapi.MessageConfig{}, err
	}
	if isCommandInput {
		return msg, nil
	}

	url, err := h.handleUrlInput(chatID, input)
	if err != nil {
		return tgbotapi.MessageConfig{}, err
	}

	if err := h.srv.CreateTask(ctx, url); err != nil {
		return tgbotapi.MessageConfig{}, err
	}

	msg = tgbotapi.NewMessage(chatID, taskWasCreatedMessage)
	h.storage.SetState(userID, types.Menu)

	return msg, nil
}

func (h *Handler) handleDeleteTaskUrlInput(ctx context.Context, input string, userID, chatID int64) (tgbotapi.MessageConfig, error) {
	msg, isCommandInput, err := h.handleCommandInput(ctx, input, userID, chatID)
	if err != nil {
		return tgbotapi.MessageConfig{}, err
	}
	if isCommandInput {
		return msg, nil
	}

	url, err := h.handleUrlInput(chatID, input)
	if err != nil {
		return tgbotapi.MessageConfig{}, err
	}

	if err := h.srv.DeleteTask(ctx, url); err != nil {
		return tgbotapi.MessageConfig{}, err
	}

	msg = tgbotapi.NewMessage(chatID, taskWasDeletedMessage)
	h.storage.SetState(userID, types.Menu)

	return msg, nil
}

func (h *Handler) handleAssignWorkerUrlInput(ctx context.Context, input string, userID, chatID int64) (tgbotapi.MessageConfig, error) {
	msg, isCommandInput, err := h.handleCommandInput(ctx, input, userID, chatID)
	if err != nil {
		return tgbotapi.MessageConfig{}, err
	}
	if isCommandInput {
		return msg, nil
	}

	url, err := h.handleUrlInput(chatID, input)
	if err != nil {
		return tgbotapi.MessageConfig{}, err
	}

	msg = tgbotapi.NewMessage(chatID, enterUserIdMessage)
	h.storage.SetStateWithData(userID, types.AssignWorkerIdInput, types.UrlData, url)

	return msg, nil
}

func (h *Handler) handleAssignWorkerIdInput(ctx context.Context, input string, userID, chatID int64) (tgbotapi.MessageConfig, error) {
	msg, isCommandInput, err := h.handleCommandInput(ctx, input, userID, chatID)
	if err != nil {
		return tgbotapi.MessageConfig{}, err
	}
	if isCommandInput {
		return msg, nil
	}

	id, err := h.handleIdInput(chatID, input)
	if err != nil {
		return tgbotapi.MessageConfig{}, err
	}

	url, ok := h.storage.GetData(userID, types.UrlData)
	if !ok {
		return tgbotapi.MessageConfig{}, errGettingUrlFromStorage
	}

	if err := h.srv.AssignUser(ctx, url.(string), id); err != nil {
		return tgbotapi.MessageConfig{}, err
	}

	msg = tgbotapi.NewMessage(chatID, workerWasAssignedMessage)
	h.storage.SetState(userID, types.Menu)

	return msg, nil
}

func (h *Handler) handleCloseTaskUrlInput(ctx context.Context, input string, userID, chatID int64) (tgbotapi.MessageConfig, error) {
	msg, isCommandInput, err := h.handleCommandInput(ctx, input, userID, chatID)
	if err != nil {
		return tgbotapi.MessageConfig{}, err
	}
	if isCommandInput {
		return msg, nil
	}

	url, err := h.handleUrlInput(chatID, input)
	if err != nil {
		return tgbotapi.MessageConfig{}, err
	}

	if err := h.srv.CloseTask(ctx, url); err != nil {
		return tgbotapi.MessageConfig{}, err
	}

	msg = tgbotapi.NewMessage(chatID, taskWasClosedMessage)
	h.storage.SetState(userID, types.Menu)

	return msg, nil
}

func (h *Handler) handleCommandInput(ctx context.Context, input string, userID, chatID int64) (tgbotapi.MessageConfig, isCommandInput, error) {
	var msg tgbotapi.MessageConfig

	switch input {
	case startCommand:
		if err := h.srv.SignUp(ctx, userID); err != nil {
			return tgbotapi.MessageConfig{}, true, err
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
			return tgbotapi.MessageConfig{}, true, err
		}
		if tasks == nil {
			msg = tgbotapi.NewMessage(chatID, emptyTasksListMessage)
		} else {
			msg = tgbotapi.NewMessage(chatID, currentOpenTasksMessage)
			msg.ReplyMarkup = NewTasksKeyboard(tasks)
		}

	default:
		return tgbotapi.MessageConfig{}, false, nil
	}

	return msg, true, nil
}

func (h *Handler) handleError(err error, chatID int64) tgbotapi.MessageConfig {
	log.Println(err.Error())
	return tgbotapi.NewMessage(chatID, err.Error())
}

func (h *Handler) handleUrlInput(chatID int64, input string) (string, error) {
	_, err := url.ParseRequestURI(input)
	if err != nil {
		return "", errInvalidUrlInput
	}

	return input, nil
}

func (h *Handler) handleIdInput(chatID int64, input string) (int64, error) {
	id, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return 0, errInvalidUserIdInput
	}

	return id, nil
}

func getUserSignedUpMessage(userID int64) string {
	return fmt.Sprintf(userSignedUpMessageFmt, userID)
}
