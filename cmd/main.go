package main

import (
	"log"

	"github.com/smokinjoints/crypto-price-bot/pkg/coincap"
	"github.com/smokinjoints/crypto-price-bot/pkg/telegram"
)

func main() {
	cfg, err := coincap.InitConfig()
	if err != nil {
		log.Println(err)
	}

	client := coincap.NewCoincapClient()

	bot, err := telegram.NewBot(cfg)
	if err != nil {
		log.Fatal(err)
	}

	telegram.Start(bot, *cfg, *client)
}
