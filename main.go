package main

import (
	"log"

	"tg-openim/config"
	"tg-openim/controller"
	"tg-openim/router"
	"tg-openim/service"
)

func main() {

	// 配置
	config.InitConfig()

	// 创建 service 实例, Telegram Service
	telegramService, err := service.NewTelegramService(
		config.App.BotToken,
		config.App.WebhookUrl,
	)

	if err != nil {
		panic(err)
	}

	// OpenIM Service
	openIMService := service.NewOpenIMService()	

	// OpenAI Service
	openAIService := service.NewOpenAIService()

	// 获取 OpenIM token
	err = openIMService.RefreshAdminToken()

	if err != nil {
		panic(err)
	}

	// 创建 Controller 对象，并注入依赖
	tgController := controller.NewTgController(
		openIMService,
		telegramService,
		openAIService,
	)

	openIMController := controller.NewOpenIMController(
		telegramService,
	)

	// Router
	r := router.InitRouter(
		tgController,
		openIMController,
	)

	log.Println("服务启动:", config.App.Port)

	err = r.Run(":" + config.App.Port)

	if err != nil {
		panic(err)
	}
}