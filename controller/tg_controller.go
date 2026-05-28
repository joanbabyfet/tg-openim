package controller

import (
	"fmt"
	"log"
	"strings"

	"tg-openim/cache"
	"tg-openim/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/gin-gonic/gin"
)

type TgController struct {
	OpenIMService *service.OpenIMService
	TelegramService *service.TelegramService
	OpenAIService *service.OpenAIService
}

//构造函数
func NewTgController(
	openim *service.OpenIMService,
	tg *service.TelegramService,
	openai *service.OpenAIService,
) *TgController {

	return &TgController{
		OpenIMService:   openim,
		TelegramService: tg,
		OpenAIService: openai,
	}
}

//把函数改成 method, Telegram webhook
func (c *TgController) TgWebhook(ctx *gin.Context) {

	var update tgbotapi.Update

	if err := ctx.ShouldBindJSON(&update); err != nil {

		ctx.JSON(400, gin.H{
			"msg": err.Error(),
		})

		return
	}

	//InlineKeyboard 按钮点击
	if update.CallbackQuery != nil {
		c.HandleCallback(update)

		ctx.JSON(200, gin.H{
			"ok": true,
		})

		return
	}

	//普通消息
	if update.Message != nil {
		c.HandleMessage(update)
	}

	ctx.JSON(200, gin.H{
		"ok": true,
	})
}

func (c *TgController) HandleMessage(update tgbotapi.Update) {
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
	userID := c.OpenIMService.EnsureTGUser(
		tgUser.ID,
		tgUser.UserName,
		tgUser.FirstName,
	)

	// TG ↔ OpenIM 映射
	cache.TgUserMap[userID] = chatID

	//菜单命令
	switch text {

	case "/start":
		c.TelegramService.SendTelegramMessage(chatID, "欢迎使用系统", c.TelegramService.MainMenu())
		return
	case "/menu":
		c.TelegramService.SendTelegramMessage(chatID, "请选择菜单", c.TelegramService.MainMenu())
		return
	case "/help":
		helpText := `
可用命令：

/start - 启动机器人
/menu  - 显示菜单
/help  - 帮助说明

功能说明：

• 充值
• 提现
• 钱包查询
• AI客服
• 人工客服
`
		c.TelegramService.SendTelegramMessage(chatID, helpText, nil)
		return
	}

	//异步转发 AI + OpenIM（避免阻塞 webhook）
	go func(uid string, msg string, cid int64) {

		log.Println("开始AI处理:", uid)

		reply, err := c.OpenAIService.ChatAI(msg)

		// AI失败 → 转人工
		if err != nil {

			log.Println("AI失败, 转OpenIM:", err)

			err = c.OpenIMService.SendToOpenIM(uid, msg)

			if err != nil {
				log.Println("OpenIM转发失败:", err)
			} else {
				log.Println("已转OpenIM人工客服")
			}

			return
		}

		// 回复TG
		err = c.TelegramService.SendTelegramMessage(cid, reply, nil)

		if err != nil {
			log.Println("TG发送失败:", err)
		}

		// 转人工
		if strings.Contains(reply, "转人工") {

			log.Println("AI触发转人工")

			err = c.OpenIMService.SendToOpenIM(uid, msg)

			if err != nil {
				log.Println("转OpenIM失败:", err)
			} else {
				log.Println("已转OpenIM客服")
			}
		}

	}(userID, text, chatID)
}

func (c *TgController) HandleCallback(update tgbotapi.Update) {
	callback := update.CallbackQuery

	data := callback.Data

	chatID := callback.Message.Chat.ID

	switch data {

	// 充值
	case "deposit":
		c.TelegramService.SendTelegramMessage(chatID, "请输入充值金额", nil)
	// 提现
	case "withdraw":
		c.TelegramService.SendTelegramMessage(chatID, "请输入提现金额", nil)
	// 钱包
	case "wallet":
		msg := fmt.Sprintf("钱包余额: %.2f USDT", 100.50)
		c.TelegramService.SendTelegramMessage(chatID, msg, nil)
	// 客服
	case "support":
		c.TelegramService.SendTelegramMessage(chatID, "客服处理中...", nil)
	}

	// 告诉 Telegram callback 已处理
	callbackConfig := tgbotapi.NewCallback(
		callback.ID,
		"操作成功",
	)

	c.TelegramService.Bot.Request(callbackConfig)
}