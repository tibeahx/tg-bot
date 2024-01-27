package telegram

import (
	"encoding/json"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/smokinjoints/crypto-price-bot/pkg/coincap"
	"github.com/smokinjoints/crypto-price-bot/pkg/models"
)

func handleStart(b *Bot, msg *tgbotapi.Message) {
	responseMessage := "Available assets to check: sol, eth, btc"
	b.Send(tgbotapi.NewMessage(msg.Chat.ID, responseMessage))
}

func handleAsset(asset models.Asset, msg *tgbotapi.Message, cfg coincap.Config, client coincap.CoincapClient, bot *Bot) {
	response := &models.Response{}

	resp, err := coincap.GetAssetPrice(client, asset, cfg)
	if err != nil {
		log.Printf("error getting %s price: %v", asset.Id, err)
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("error getting %s price", asset.Id)))
		return
	}

	if err := json.Unmarshal(resp, response); err != nil {
		log.Printf("failed to unmarshal response into struct: %v", err)
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "failed to unmarshal response into struct"))
		return
	}

	if len(response.Data) > 0 {
		priceUsd := response.Data[0].PriceUsd
		responseMessage := fmt.Sprintf("%s price: %s", asset.Id, priceUsd)
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, responseMessage))
	}

}
