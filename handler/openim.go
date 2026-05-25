package handler

import (
	"encoding/json"
	"io"
	"log"

	"tg-openim/model"
	"tg-openim/service"

	"github.com/gin-gonic/gin"
)

type OpenIMMessage struct {
	SendID      string `json:"sendID"`
	RecvID      string `json:"recvID"`
	Content     string `json:"content"`
	ContentType int    `json:"contentType"`
}

type TextContent struct {
	Content string `json:"content"`
}

//OpenIM 回调
func OpenIMCallback(c *gin.Context) {

	log.Println("======== OpenIM Callback ========")

	body, _ := io.ReadAll(c.Request.Body)

	log.Println(string(body))

	var msg OpenIMMessage

	if err := json.Unmarshal(body, &msg); err != nil {

		c.JSON(400, gin.H{
			"err": err.Error(),
		})

		return
	}

	// 只处理文本消息
	if msg.ContentType != 101 {

		c.JSON(200, gin.H{
			"errCode": 0,
			"errMsg":  "",
		})

		return
	}

	var text TextContent

	json.Unmarshal([]byte(msg.Content), &text)
	
	chatID, ok := model.TgUserMap[msg.RecvID]
	if ok {
		service.SendTelegramMessage(
			chatID,
			"客服回复: "+text.Content,
		)
	}

	c.JSON(200, gin.H{
		"errCode": 0,
		"errMsg":  "",
	})
}