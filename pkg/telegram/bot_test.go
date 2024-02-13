package telegram

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestGetTelegramBot(t *testing.T) {
	apiKey := "mockAPIKey"

	bot := GetTelegramBot(apiKey)

	if bot == nil {
		t.Errorf("function call didn't return the bot")
	}
}

func TestSend(t *testing.T) {
	msg := tgbotapi.NewMessage(123, "test message")

	err := Send(msg)

	if err != nil {
		t.Errorf("Send() returned an error: %v", err)
	}
}
