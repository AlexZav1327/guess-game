package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	tgbot "github.com/AlexZav1327/guess-game/internal/tg-bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

func main() {
	botToken := os.Getenv("BOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "main",
			"function": "NewBotAPI",
			"error":    err,
		}).Fatal("BotAPI instance creation error")

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
