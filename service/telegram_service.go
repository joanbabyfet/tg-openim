package service

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramService struct {
	Bot *tgbotapi.BotAPI
}

// 构造函数
func NewTelegramService(
	token string,
	webhookURL string,
) (*TelegramService, error) {

	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		return nil, err
	}

	bot.Debug = true

	// webhook
	webhook, err := tgbotapi.NewWebhook(
		webhookURL,
	)

	if err != nil {
		return nil, err
	}

	_, err = bot.Request(webhook)

	if err != nil {
		return nil, err
	}

	info, err := bot.GetWebhookInfo()

	if err != nil {
		return nil, err
	}

	log.Printf("Webhook: %+v\n", info)

	log.Println("Telegram Bot 已启动")

	return &TelegramService{
		Bot: bot,
	}, nil
}

//发消息
func (s *TelegramService) SendTelegramMessage(chatID int64, text string, replyMarkup interface{}) error {
	msg := tgbotapi.NewMessage(chatID, text)

	// 设置按钮
	if replyMarkup != nil {
		msg.ReplyMarkup = replyMarkup
	}
	
	_, err := s.Bot.Send(msg)

	if err != nil {
		log.Println("发送TG失败:", err)
		return err
	}

	return nil
}

//主菜单
func (s *TelegramService) MainMenu() tgbotapi.InlineKeyboardMarkup {

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