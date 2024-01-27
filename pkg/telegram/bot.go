package telegram

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/smokinjoints/crypto-price-bot/pkg/coincap"
	"github.com/smokinjoints/crypto-price-bot/pkg/models"
)

type Bot struct {
	bot *tgbotapi.BotAPI
}

func NewBot(cfg *coincap.Config) (*Bot, error) {
	botAPI, err := tgbotapi.NewBotAPI(cfg.BotAPIkey)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	botAPI.Debug = true
	return &Bot{bot: botAPI}, nil
}

func (b *Bot) Send(msg tgbotapi.Chattable) error {
	_, err := b.bot.Send(msg)
	return err
}

func Start(b *Bot, cfg coincap.Config, client coincap.CoincapClient) {
	updCfg := tgbotapi.NewUpdate(0)
	updCfg.Timeout = 60

	for {
		updates := b.bot.GetUpdatesChan(updCfg)

		for update := range updates {
			if update.Message != nil && update.Message.IsCommand() {
				var asset models.Asset
				switch update.Message.Command() {
				case "/start":
					handleStart(b, update.Message)

				case "/btc":
					asset = models.Asset{Id: "bitcoin"}

				case "/eth":
					asset = models.Asset{Id: "ethereum"}

				case "/sol":
					asset = models.Asset{Id: "solana"}

				default:
					asset = models.Asset{Id: update.Message.Command()}
				}

				handleAsset(asset, update.Message, cfg, client, b)
			}
		}
	}
}
