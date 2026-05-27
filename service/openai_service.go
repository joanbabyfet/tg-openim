package service

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"tg-openim/config"
)

func ChatAI(text string) (string, error) {
	log.Println("===== ChatAI Request Start =====")
	log.Println("User Input:", text)
	
	body := map[string]interface{}{
		"model": "gpt-4.1-mini", //用哪个 AI 模型
		"messages": []map[string]string{
			{
				"role":    "system",	//你是一个客服助手
				"content": "你是一个客服助手",
			},
			{
				"role": "user", //就是用户说的话
				"content": text,
			},
		},
	}

	bs, _ := json.Marshal(body)

	req, _ := http.NewRequest(
		"POST",
		"https://api.openai.com/v1/chat/completions",
		bytes.NewBuffer(bs),
	)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.App.OpenAIKey)

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)
	//AI返回文字内容
	content := result["choices"].([]interface{})[0].
		(map[string]interface{})["message"].
		(map[string]interface{})["content"].(string)

	log.Println("AI Response:", content)
	log.Println("===== ChatAI Request End =====")

	return content, nil
}