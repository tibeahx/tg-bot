package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tibeahx/tg-bot/pkg/coincap"
	"github.com/tibeahx/tg-bot/pkg/models"
)

var bot *tgbotapi.BotAPI

func GetTelegramBot(apiKey string) *tgbotapi.BotAPI {
	botAPI, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		log.Fatal(err)
	}

	botAPI.Debug = true
	bot = botAPI

	return bot
}

func Send(msg tgbotapi.Chattable) error {
	_, err := bot.Send(msg)
	return err
}

func Start(client coincap.CoincapClient) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	for {
		updates := bot.GetUpdatesChan(updateConfig)

		for update := range updates {
			if update.Message != nil && update.Message.IsCommand() {
				var asset models.Asset
				switch update.Message.Text {
				case "/start":
					handleStart(update.Message)
					continue

				case "/btc":
					asset = models.Asset{Name: "bitcoin"}
					handleAsset(asset, update.Message, client, bot)
					continue

				case "/eth":
					asset = models.Asset{Name: "ethereum"}
					handleAsset(asset, update.Message, client, bot)
					continue

				case "/sol":
					asset = models.Asset{Name: "solana"}
					handleAsset(asset, update.Message, client, bot)
					continue
				}
			}
		}
	}
}
