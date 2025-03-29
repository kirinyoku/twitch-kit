package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kirinyoku/twitch-kit/internal/bot"
	"github.com/kirinyoku/twitch-kit/internal/fetcher"
	"github.com/kirinyoku/twitch-kit/pkg/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	botAPI, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Printf("failed to initialize bot: %v", err)
		return
	}

	fetcher := fetcher.NewFetcher()

	tgBot := bot.New(botAPI, fetcher)
	tgBot.RegisterCommand("start", bot.ViewCmdStart())

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := tgBot.Start(ctx); err != nil {
		log.Printf("failed to start bot: %v", err)
		return
	}
}
