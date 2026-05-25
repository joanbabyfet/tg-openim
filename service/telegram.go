package service

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Bot *tgbotapi.BotAPI

func InitTelegram(token string, webhookURL string) error {

	var err error

	Bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	// 开启 debug（可选）
	Bot.Debug = true

	// 创建 webhook
	webhook, err := tgbotapi.NewWebhook(webhookURL)
	if err != nil {
		return err
	}

	// 设置 webhook
	_, err = Bot.Request(webhook)
	if err != nil {
		return err
	}

	// 查看 webhook 信息（可选）
	info, err := Bot.GetWebhookInfo()
	if err != nil {
		return err
	}

	log.Printf("Webhook: %+v\n", info)
	log.Println("Telegram Bot 已启动")

	return nil
}

func SendTelegramMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)

	_, err := Bot.Send(msg)

	if err != nil {
		log.Println("发送TG失败:", err)
		return
	}
}