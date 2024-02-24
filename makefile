.PHONY
.SILENT

build:
	go build -o ./.bin/tg-bot cmd/tg-bot/main.go

run: build
	./.bin/bot

test:
	go test ./...
