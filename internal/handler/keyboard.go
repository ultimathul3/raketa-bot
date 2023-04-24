package handler

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vanyaio/raketa-bot/internal/types"
)

const (
	startCommand        = "/start"
	createTaskCommand   = "Create task"
	deleteTaskCommand   = "Delete task"
	assignWorkerCommand = "Assign worker"
	closeTaskCommand    = "Close task"
	getOpenTasksCommand = "Get open tasks"
)

var menuKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(createTaskCommand+" ➕"),
		tgbotapi.NewKeyboardButton(deleteTaskCommand+" ➖"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(getOpenTasksCommand+" 📃"),
		tgbotapi.NewKeyboardButton(closeTaskCommand+" ✔"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(assignWorkerCommand+" 👨‍🔧"),
	),
)

func NewTasksKeyboard(tasks []types.Task) tgbotapi.InlineKeyboardMarkup {
	var buttons [][]tgbotapi.InlineKeyboardButton

	for i, task := range tasks {
		buttons = append(buttons,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL(
					fmt.Sprintf("Task %d", i+1),
					task.Url,
				),
			),
		)
	}

	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}
