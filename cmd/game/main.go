package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	tgbot "github.com/AlexZav1327/guess-game/internal/tg-bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	var botToken = os.Getenv("BOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		fmt.Println("BotAPI instance creation error:", err)
		return
	}

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	Updates := bot.GetUpdatesChan(updateConfig)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	defer cancel()

	gb := tgbot.GuessBot{Updates: Updates, Bot: *bot}
	gb.Run(ctx)
}
