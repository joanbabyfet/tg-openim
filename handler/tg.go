package handler

import (
	"fmt"
	"log"
	"strings"

	"tg-openim/model"
	"tg-openim/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/gin-gonic/gin"
)

//Telegram webhook
func TgWebhook(c *gin.Context) {

	var update tgbotapi.Update

	if err := c.ShouldBindJSON(&update); err != nil {

		c.JSON(400, gin.H{
			"msg": err.Error(),
		})

		return
	}

	//InlineKeyboard 按钮点击
	if update.CallbackQuery != nil {
		HandleCallback(update)

		c.JSON(200, gin.H{
			"ok": true,
		})

		return
	}

	//普通消息
	if update.Message != nil {
		HandleMessage(update)
	}

	c.JSON(200, gin.H{
		"ok": true,
	})
}

func HandleMessage(update tgbotapi.Update) {
	text := update.Message.Text
	chatID := update.Message.Chat.ID
	tgUser := update.Message.From

	// TG显示名
	userName := tgUser.UserName
	if userName == "" {
		userName = tgUser.FirstName
	}

	log.Println("TG收到消息:", text, userName)

	// 自动注册 OpenIM 用户
	userID := service.EnsureTGUser(
		tgUser.ID,
		tgUser.UserName,
		tgUser.FirstName,
	)

	// TG ↔ OpenIM 映射
	model.TgUserMap[userID] = chatID

	//菜单命令
	switch text {

	case "/start":
		service.SendTelegramMessage(chatID, "欢迎使用系统", service.MainMenu())
	default:
		service.SendTelegramMessage(chatID, "请选择菜单", service.MainMenu())
	}

	//异步转发 AI + OpenIM（避免阻塞 webhook）
	go func(uid string, msg string, cid int64) {

		log.Println("开始AI处理:", uid)

		reply, err := service.ChatAI(msg)

		// AI失败 → 转人工
		if err != nil {

			log.Println("AI失败, 转OpenIM:", err)

			err = service.SendToOpenIM(uid, msg)

			if err != nil {
				log.Println("OpenIM转发失败:", err)
			} else {
				log.Println("已转OpenIM人工客服")
			}

			return
		}

		// 回复TG
		err = service.SendTelegramMessage(cid, reply, nil)

		if err != nil {
			log.Println("TG发送失败:", err)
		}

		// 转人工
		if strings.Contains(reply, "转人工") {

			log.Println("AI触发转人工")

			err = service.SendToOpenIM(uid, msg)

			if err != nil {
				log.Println("转OpenIM失败:", err)
			} else {
				log.Println("已转OpenIM客服")
			}
		}

	}(userID, text, chatID)
}

func HandleCallback(update tgbotapi.Update) {
	callback := update.CallbackQuery

	data := callback.Data

	chatID := callback.Message.Chat.ID

	switch data {

	// 充值
	case "deposit":
		service.SendTelegramMessage(chatID, "请输入充值金额", nil)
	// 提现
	case "withdraw":
		service.SendTelegramMessage(chatID, "请输入提现金额", nil)
	// 钱包
	case "wallet":
		msg := fmt.Sprintf("钱包余额: %.2f USDT", 100.50)
		service.SendTelegramMessage(chatID, msg, nil)
	// 客服
	case "support":
		service.SendTelegramMessage(chatID, "客服处理中...", nil)
	}

	// 告诉 Telegram callback 已处理
	callbackConfig := tgbotapi.NewCallback(
		callback.ID,
		"操作成功",
	)

	service.Bot.Request(callbackConfig)
}