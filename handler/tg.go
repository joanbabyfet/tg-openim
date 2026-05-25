package handler

import (
	"log"

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
			"error": err.Error(),
		})

		return
	}

	if update.Message == nil {

		c.JSON(200, gin.H{
			"ok": true,
		})

		return
	}

	text := update.Message.Text
	chatID := update.Message.Chat.ID
	tgUser := update.Message.From
	
	// 获取tg显示名
	userName := tgUser.UserName
	if userName == "" {
		userName = tgUser.FirstName
	}
	log.Println("TG收到消息:", text)

	// 自动注册 OpenIM 用户
	userID := service.EnsureTGUser(
		tgUser.ID, //识别“谁发的”
		tgUser.UserName,
		tgUser.FirstName,
	)

	// 保存 TG ↔ OpenIM 映射
	model.TgUserMap[userID] = chatID

	// 异步转发 OpenIM
	go func() {

		err := service.SendToOpenIM(userID, text)

		if err != nil {
			log.Println("转发OpenIM失败:", err)
		}
	}()

	c.JSON(200, gin.H{
		"ok": true,
	})
}