package telegram

import (
	"encoding/json"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/smokinjoints/crypto-price-bot/pkg/coincap"
	"github.com/smokinjoints/crypto-price-bot/pkg/models"
)

func handleStart(msg *tgbotapi.Message) error {
	responseMessage := "Available assets to check: sol, eth, btc"
	bot.Send(tgbotapi.NewMessage(msg.Chat.ID, responseMessage))

	return nil
}

func handleAsset(asset models.Asset, msg *tgbotapi.Message, client coincap.CoincapClient, bot *tgbotapi.BotAPI) error {
	resp, err := coincap.GetAssetPrice(client, asset)
	if err != nil {
		log.Printf("error getting %s price: %v", asset.Name, err)
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("error getting %s price", asset.Name)))
		return nil
	}

	response := &models.Response{}

	if err := json.Unmarshal(resp, &response); err != nil {
		log.Printf("failed to unmarshal response into struct: %v", err)
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "failed to unmarshal response into struct"))
		return nil
	}

	if response.Data.PriceUsd != "" {
		priceUsd := response.Data.PriceUsd
		responseMessage := fmt.Sprintf("%s price: %s", asset.Name, priceUsd)
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, responseMessage))
	}

	return nil
}
