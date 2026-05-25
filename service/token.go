package service

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"tg-openim/config"
)

var AdminToken string

//获取imAdmin管理员token
func RefreshAdminToken() error {

	body := map[string]interface{}{
		"userID": config.App.OpenIMAdmin,
		"secret": config.App.OpenIMSecret,
	}

	bs, _ := json.Marshal(body)

	//获取imAdmin管理员token
	req, _ := http.NewRequest(
		"POST",
		config.App.OpenIMAPI+"/auth/get_admin_token",
		bytes.NewBuffer(bs),
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("operationID", strconv.FormatInt(time.Now().UnixMilli(), 10))

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	log.Println("token返回:", result)

	data := result["data"].(map[string]interface{})
	AdminToken = data["token"].(string)

	log.Println("OpenIM token 获取成功")

	return nil
}