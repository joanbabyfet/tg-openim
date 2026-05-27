package main

import (
	"log"

	"tg-openim/config"
	"tg-openim/router"
	"tg-openim/service"
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

	// 获取 OpenIM 管理员 token
	err = service.RefreshAdminToken()
	if err != nil {
		panic(err)
	}

	// 路由
	r := router.InitRouter()

	log.Println("服务启动:", config.App.Port)

	r.Run(":" + config.App.Port)
}