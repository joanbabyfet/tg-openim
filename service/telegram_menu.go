// internal/service/telegram_menu.go

package service

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func MainMenu() tgbotapi.InlineKeyboardMarkup {

    return tgbotapi.NewInlineKeyboardMarkup(

		tgbotapi.NewInlineKeyboardRow(

			tgbotapi.NewInlineKeyboardButtonData(
				"💰 充值",
				"deposit",
			),

			tgbotapi.NewInlineKeyboardButtonData(
				"💸 提现",
				"withdraw",
			),
		),

		tgbotapi.NewInlineKeyboardRow(

			tgbotapi.NewInlineKeyboardButtonData(
				"👛 钱包",
				"wallet",
			),

			tgbotapi.NewInlineKeyboardButtonData(
				"🎧 客服",
				"support",
			),
		),
	)
}