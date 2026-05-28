package router

import (
	"tg-openim/controller"

	"github.com/gin-gonic/gin"
)

//改为 注入 controller 实例
func InitRouter(
	tgController *controller.TgController,
	openIMController *controller.OpenIMController,
) *gin.Engine {

	r := gin.Default()

	// TG webhook
	r.POST("/webhook", tgController.TgWebhook)

	// OpenIM callback
	r.POST("/openim/callback/callbackAfterSendSingleMsgCommand", openIMController.OpenIMCallback)

	return r
}