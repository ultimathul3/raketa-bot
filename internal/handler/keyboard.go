package handler

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	raketapb "github.com/vanyaio/raketa-backend/proto"
)

var menuKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Create task ➕"),
		tgbotapi.NewKeyboardButton("Delete task ➖"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Get open tasks 📃"),
		tgbotapi.NewKeyboardButton("Close task ✔"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Assign worker 👨‍🔧"),
	),
)

func NewTasksKeyboard(tasks []*raketapb.Task) tgbotapi.InlineKeyboardMarkup {
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
