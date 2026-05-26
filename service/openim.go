package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"tg-openim/config"
	"tg-openim/model"
)

func SendToOpenIM(userID string, text string) error {

	body := map[string]interface{}{
		"sendID": userID, //tg_123456 用户
		"recvID": config.App.OpenIMCustomerService, // OpenIM接收用户(客服)
		"groupID": "",
		"senderPlatformID": 1,
		"sessionType": 1,
		"contentType": 101,
		"content": map[string]interface{}{
			"content": text,
		},
	}

	bs, _ := json.Marshal(body)

	req, _ := http.NewRequest(
		"POST",
		config.App.OpenIMAPI+"/msg/send_msg",
		bytes.NewBuffer(bs),
	)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("token", AdminToken)
	req.Header.Set("operationID", strconv.FormatInt(time.Now().UnixMilli(), 10))

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	log.Println("OpenIM响应:", result)

	return nil
}

func RegisterOpenIMUser(userID string, nickname string) error {

	body := model.OpenIMRegisterReq{
		Users: []model.OpenIMUser{
			{
				UserID:   userID,
				Nickname: nickname,
			},
		},
	}

	bs, _ := json.Marshal(body)

	req, _ := http.NewRequest(
		"POST",
		config.App.OpenIMAPI+"/user/user_register",
		bytes.NewBuffer(bs),
	)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("token", AdminToken)
	req.Header.Set("operationID", strconv.FormatInt(time.Now().UnixMilli(), 10))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	log.Println("注册用户:", result)

	return nil
}

func EnsureTGUser(tgID int64, username string, firstName string) string {
	//不要用tg username当当 OpenIM 用户ID, 因tg username用户可修改
	userID := fmt.Sprintf("tg_%d", tgID)

	nickname := username

	if nickname == "" {
		nickname = firstName
	}

	_ = RegisterOpenIMUser(userID, nickname)

	return userID
}