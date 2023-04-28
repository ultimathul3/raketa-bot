package handler

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"log"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vanyaio/raketa-bot/internal/types"
)

type storage interface {
	GetState(userID int64) (types.State, bool)
	GetData(userID int64, key string) (any, bool)
	SetState(userID int64, state types.State)
	SetStateWithData(userID int64, state types.State, key string, data any)
}

type service interface {
	SignUp(ctx context.Context, id int64, username string) error
	CreateTask(ctx context.Context, url string) error
	SetTaskPrice(ctx context.Context, url string, price uint64) error
	DeleteTask(ctx context.Context, url string) error
	AssignUser(ctx context.Context, url, username string) error
	CloseTask(ctx context.Context, url string) error
	GetUnassignTasks(ctx context.Context) ([]types.Task, error)
}

type isCommandInput bool

type Handler struct {
	srv     service
	bot     *tg.BotAPI
	storage storage
}

func NewHandler(srv service, bot *tg.BotAPI, storage storage) *Handler {
	return &Handler{
		srv:     srv,
		bot:     bot,
		storage: storage,
	}
}

func (h *Handler) HandleUpdates(ctx context.Context, config tg.UpdateConfig) {
	updates := h.bot.GetUpdatesChan(config)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		inputMessage := update.Message
		message := tg.MessageConfig{}
		var err error

		state, ok := h.storage.GetState(inputMessage.From.ID)
		if !ok {
			state = types.Menu
			h.storage.SetState(inputMessage.From.ID, state)
		}

		switch state {
		case types.Menu:
			message, _, err = h.handleCommandInput(ctx, inputMessage)
			if err != nil {
				message = h.handleError(err, inputMessage.Chat.ID)
			}

		case types.CreateTaskUrlInput:
			message, err = h.handleCreateTaskUrlInput(ctx, inputMessage)
			if err != nil {
				message = h.handleError(err, inputMessage.Chat.ID)
			}

		case types.CreateTaskPriceInput:
			message, err = h.handleCreateTaskPriceInput(ctx, inputMessage)
			if err != nil {
				message = h.handleError(err, inputMessage.Chat.ID)
			}

		case types.DeleteTaskUrlInput:
			message, err = h.handleDeleteTaskUrlInput(ctx, inputMessage)
			if err != nil {
				message = h.handleError(err, inputMessage.Chat.ID)
			}

		case types.AssignWorkerUrlInput:
			message, err = h.handleAssignWorkerUrlInput(ctx, inputMessage)
			if err != nil {
				message = h.handleError(err, inputMessage.Chat.ID)
			}

		case types.AssignWorkerUsernameInput:
			message, err = h.handleAssignWorkerUsernameInput(ctx, inputMessage)
			if err != nil {
				message = h.handleError(err, inputMessage.Chat.ID)
			}

		case types.CloseTaskUrlInput:
			message, err = h.handleCloseTaskUrlInput(ctx, inputMessage)
			if err != nil {
				message = h.handleError(err, inputMessage.Chat.ID)
			}
		}

		if message.Text != "" {
			h.bot.Send(message)
		}
	}
}

func (h *Handler) handleCreateTaskUrlInput(ctx context.Context, input *tg.Message) (tg.MessageConfig, error) {
	msg, isCommandInput, err := h.handleCommandInput(ctx, input)
	if err != nil {
		return tg.MessageConfig{}, err
	}
	if isCommandInput {
		return msg, nil
	}

	url, err := h.handleUrlInput(input.Text)
	if err != nil {
		return tg.MessageConfig{}, err
	}

	msg = tg.NewMessage(input.Chat.ID, enterTaskPriceMessage)
	h.storage.SetStateWithData(input.From.ID, types.CreateTaskPriceInput, types.UrlDataKey, url)

	return msg, nil
}

func (h *Handler) handleCreateTaskPriceInput(ctx context.Context, input *tg.Message) (tg.MessageConfig, error) {
	msg, isCommandInput, err := h.handleCommandInput(ctx, input)
	if err != nil {
		return tg.MessageConfig{}, err
	}
	if isCommandInput {
		return msg, nil
	}

	price, err := h.handlePriceInput(input.Text)
	if err != nil {
		return tg.MessageConfig{}, err
	}

	url, ok := h.storage.GetData(input.From.ID, types.UrlDataKey)
	if !ok {
		return tg.MessageConfig{}, errGettingUrlFromStorage
	}
	if err := h.srv.CreateTask(ctx, url.(string)); err != nil {
		return tg.MessageConfig{}, err
	}

	if err := h.srv.SetTaskPrice(ctx, url.(string), price); err != nil {
		return tg.MessageConfig{}, err
	}

	msg = tg.NewMessage(input.Chat.ID, taskWasCreatedMessage)
	h.storage.SetState(input.From.ID, types.Menu)

	return msg, nil
}

func (h *Handler) handleDeleteTaskUrlInput(ctx context.Context, input *tg.Message) (tg.MessageConfig, error) {
	msg, isCommandInput, err := h.handleCommandInput(ctx, input)
	if err != nil {
		return tg.MessageConfig{}, err
	}
	if isCommandInput {
		return msg, nil
	}

	url, err := h.handleUrlInput(input.Text)
	if err != nil {
		return tg.MessageConfig{}, err
	}

	if err := h.srv.DeleteTask(ctx, url); err != nil {
		return tg.MessageConfig{}, err
	}

	msg = tg.NewMessage(input.Chat.ID, taskWasDeletedMessage)
	h.storage.SetState(input.From.ID, types.Menu)

	return msg, nil
}

func (h *Handler) handleAssignWorkerUrlInput(ctx context.Context, input *tg.Message) (tg.MessageConfig, error) {
	msg, isCommandInput, err := h.handleCommandInput(ctx, input)
	if err != nil {
		return tg.MessageConfig{}, err
	}
	if isCommandInput {
		return msg, nil
	}

	url, err := h.handleUrlInput(input.Text)
	if err != nil {
		return tg.MessageConfig{}, err
	}

	msg = tg.NewMessage(input.Chat.ID, enterUsernameMessage)
	h.storage.SetStateWithData(input.From.ID, types.AssignWorkerUsernameInput, types.UrlDataKey, url)

	return msg, nil
}

func (h *Handler) handleAssignWorkerUsernameInput(ctx context.Context, input *tg.Message) (tg.MessageConfig, error) {
	msg, isCommandInput, err := h.handleCommandInput(ctx, input)
	if err != nil {
		return tg.MessageConfig{}, err
	}
	if isCommandInput {
		return msg, nil
	}

	username := input.Text

	url, ok := h.storage.GetData(input.From.ID, types.UrlDataKey)
	if !ok {
		return tg.MessageConfig{}, errGettingUrlFromStorage
	}

	if err := h.srv.AssignUser(ctx, url.(string), username); err != nil {
		return tg.MessageConfig{}, err
	}

	msg = tg.NewMessage(input.Chat.ID, workerWasAssignedMessage)
	h.storage.SetState(input.From.ID, types.Menu)

	return msg, nil
}

func (h *Handler) handleCloseTaskUrlInput(ctx context.Context, input *tg.Message) (tg.MessageConfig, error) {
	msg, isCommandInput, err := h.handleCommandInput(ctx, input)
	if err != nil {
		return tg.MessageConfig{}, err
	}
	if isCommandInput {
		return msg, nil
	}

	url, err := h.handleUrlInput(input.Text)
	if err != nil {
		return tg.MessageConfig{}, err
	}

	if err := h.srv.CloseTask(ctx, url); err != nil {
		return tg.MessageConfig{}, err
	}

	msg = tg.NewMessage(input.Chat.ID, taskWasClosedMessage)
	h.storage.SetState(input.From.ID, types.Menu)

	return msg, nil
}

func (h *Handler) handleCommandInput(ctx context.Context, input *tg.Message) (tg.MessageConfig, isCommandInput, error) {
	var msg tg.MessageConfig

	switch input.Text {
	case startCommand:
		if err := h.srv.SignUp(ctx, input.From.ID, input.From.UserName); err != nil {
			return tg.MessageConfig{}, true, err
		}
		msg = tg.NewMessage(input.Chat.ID, getUserSignedUpMessage(input.From.ID, input.From.UserName))
		msg.ReplyMarkup = menuKeyboard

	case createTaskCommand:
		msg = tg.NewMessage(input.Chat.ID, enterTaskUrlMessage)
		h.storage.SetState(input.From.ID, types.CreateTaskUrlInput)

	case deleteTaskCommand:
		msg = tg.NewMessage(input.Chat.ID, enterTaskUrlMessage)
		h.storage.SetState(input.From.ID, types.DeleteTaskUrlInput)

	case assignWorkerCommand:
		msg = tg.NewMessage(input.Chat.ID, enterTaskUrlMessage)
		h.storage.SetState(input.From.ID, types.AssignWorkerUrlInput)

	case closeTaskCommand:
		msg = tg.NewMessage(input.Chat.ID, enterTaskUrlMessage)
		h.storage.SetState(input.From.ID, types.CloseTaskUrlInput)

	case getOpenTasksCommand:
		tasks, err := h.srv.GetUnassignTasks(ctx)
		if err != nil {
			return tg.MessageConfig{}, true, err
		}
		if tasks == nil {
			msg = tg.NewMessage(input.Chat.ID, emptyTasksListMessage)
		} else {
			msg = tg.NewMessage(input.Chat.ID, currentOpenTasksMessage)
			msg.ReplyMarkup = NewTasksKeyboard(tasks)
		}

	default:
		return tg.MessageConfig{}, false, nil
	}

	return msg, true, nil
}

func (h *Handler) handleError(err error, chatID int64) tg.MessageConfig {
	log.Println(err.Error())
	return tg.NewMessage(chatID, err.Error())
}

func (h *Handler) handleUrlInput(input string) (string, error) {
	_, err := url.ParseRequestURI(input)
	if err != nil {
		return "", errInvalidUrlInput
	}

	return input, nil
}

func (h *Handler) handlePriceInput(input string) (uint64, error) {
	price, err := strconv.ParseUint(input, 10, 64)
	if err != nil {
		return 0, errInvalidPriceInput
	}

	return price, nil
}

func getUserSignedUpMessage(userID int64, username string) string {
	return fmt.Sprintf(userSignedUpMessageFmt, userID, username)
}
