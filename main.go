package main

import (
	"log"

	"tg-openim/config"
	"tg-openim/handler"
	"tg-openim/service"

	"github.com/gin-gonic/gin"
)

func main() {

	// 配置
	config.InitConfig()

	// 初始化TG
	err := service.InitTelegram(
		config.App.BotToken,
		config.App.WebhookUrl,
	)

	if err != nil {
		panic(err)
	}

	// 获取OpenIM 管理员token
	err = service.RefreshAdminToken()

	if err != nil {
		panic(err)
	}

	r := gin.Default()

	// TG webhook
	r.POST("/webhook", handler.TgWebhook)

	// OpenIM callback
	r.POST("/openim/callback/callbackAfterSendSingleMsgCommand", handler.OpenIMCallback)

	log.Println("服务启动:", config.App.Port)

	r.Run(":" + config.App.Port)
}