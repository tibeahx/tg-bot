package telegram

import (
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/smokinjoints/crypto-price-bot/pkg/coincap"
	"github.com/smokinjoints/crypto-price-bot/pkg/models"
)

var bot *tgbotapi.BotAPI

var once sync.Once

func GetTelegramBot(apiKey string) *tgbotapi.BotAPI {
	once.Do(func() {
		botAPI, err := tgbotapi.NewBotAPI(apiKey)
		if err != nil {
			log.Fatal(err)
		}

		botAPI.Debug = true
		bot = botAPI
	})

	return bot
}

func Send(msg tgbotapi.Chattable) error {
	_, err := bot.Send(msg)
	return err
}

func Start(cfg coincap.Config, client coincap.CoincapClient) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	for {
		updates := bot.GetUpdatesChan(updCfg)

		for update := range updates {
			if update.Message != nil && update.Message.IsCommand() {
				var asset models.Asset
				switch update.Message.Text {
				case "/start":
					handleStart(update.Message)

				case "/btc":
					asset = models.Asset{Name: "bitcoin"}
					handleAsset(asset, update.Message, client, bot)

				case "/eth":
					asset = models.Asset{Name: "ethereum"}
					handleAsset(asset, update.Message, client, bot)

				case "/sol":
					asset = models.Asset{Name: "solana"}
					handleAsset(asset, update.Message, client, bot)
				}
			}
		}
	}
}
