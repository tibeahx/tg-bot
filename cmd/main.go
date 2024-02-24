package main

import (
	"github.com/tibeahx/tg-bot/pkg/coincap"
)

func main() {
	client := coincap.NewCoincapClient()

	bot.Start(*client)
}
