package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken     string
	WebhookUrl   string
	OpenIMAPI    string
	OpenIMSecret string
	OpenIMAdmin  string
	Port         string
	OpenAIKey	 string
}

var App Config

func InitConfig() {

	_ = godotenv.Load()

	App = Config{
		BotToken:     os.Getenv("BOT_TOKEN"),
		WebhookUrl:   os.Getenv("WEBHOOK_URL"),
		OpenIMAPI:    os.Getenv("OPENIM_API"),
		OpenIMSecret: os.Getenv("OPENIM_SECRET"),
		OpenIMAdmin:  os.Getenv("OPENIM_ADMIN"),
		Port:         os.Getenv("PORT"),
		OpenAIKey:    os.Getenv("OpenAI_Key"),
	}
}