package router

import (
	"tg-openim/controller"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {

	r := gin.Default()

	// TG webhook
	r.POST("/webhook", controller.TgWebhook)

	// OpenIM callback
	r.POST("/openim/callback/callbackAfterSendSingleMsgCommand", controller.OpenIMCallback)

	return r
}