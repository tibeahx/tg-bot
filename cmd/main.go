package main

import (
	"github.com/smokinjoints/crypto-price-bot/pkg/coincap"
	"github.com/smokinjoints/crypto-price-bot/pkg/telegram"
)

func main() {
	client := coincap.NewCoincapClient()

	go telegram.Start(*client)
}
