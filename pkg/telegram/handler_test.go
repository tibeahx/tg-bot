// handler_test.go
package telegram

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/smokinjoints/crypto-price-bot/pkg/coincap"
	"github.com/smokinjoints/crypto-price-bot/pkg/models"
)

func TestHandleStart(t *testing.T) {
	msg := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: 123},
	}

	err := handleStart(msg)

	if err != nil {
		t.Errorf("got unexpected error %v", err)
	}
}

func TestHandleAsset(t *testing.T) {
	asset := models.Asset{Name: "bitcoin"}
	msg := &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 123}}
	client := coincap.CoincapClient{}
	bot := &tgbotapi.BotAPI{}

	err := handleAsset(asset, msg, client, bot)

	if err != nil {
		t.Errorf("got unexpected error: %v", err)
	}
}
