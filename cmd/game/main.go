package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/AlexZav1327/guess-game/internal/gameplay"
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
		}).Warning("BotAPI instance creation error")

		return
	}

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := bot.GetUpdatesChan(updateConfig)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	defer cancel()

	settings := gameplay.NewGameSettings(10, 1, 1000)
	game := gameplay.NewGame(10)
	guessBot := tgbot.NewBot(updates, *bot, *game, *settings)

	guessBot.Run(ctx)
}
