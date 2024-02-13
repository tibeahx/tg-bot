package main

import (
	"github.com/smokinjoints/crypto-price-bot/pkg/coincap"
	"github.com/smokinjoints/crypto-price-bot/pkg/telegram"
)

func main() {
	client := coincap.NewCoincapClient()

	telegram.Start(*client)
}
