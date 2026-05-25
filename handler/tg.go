package handler

import (
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

	// 异步转发 AI + OpenIM（避免阻塞 webhook）
	go func(uid string, msg string, cid int64) {

		log.Println("开始AI处理:", uid)

		// 1. 先走 AI
		reply, err := service.ChatAI(msg)

		// AI失败 → 直接转 OpenIM
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

		log.Println("AI回复:", reply)

		// 2. 回 Telegram
		err = service.SendTelegramMessage(cid, reply)

		if err != nil {
			log.Println("TG发送失败:", err)
		}

		// 3. 判断是否转人工
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

	c.JSON(200, gin.H{
		"ok": true,
	})
}