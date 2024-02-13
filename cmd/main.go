package main

import (
	"github.com/smokinjoints/crypto-price-bot/pkg/coincap"
	"github.com/smokinjoints/crypto-price-bot/pkg/telegram"
)

func main() {
	cfg := coincap.ReadConfig()

	client := coincap.NewCoincapClient()

	go telegram.Start(*cfg, *client)
}
