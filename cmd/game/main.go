package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	gameplay "github.com/AlexZav1327/guess-game/internal/gameplay"
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
		}).Panic("BotAPI instance creation error")

		return
	}

	debugMode := os.Getenv("DEBUG_MODE")

	bot.Debug = debugMode == "true"

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := bot.GetUpdatesChan(updateConfig)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	defer cancel()

	config := gameplay.NewDefaultConfiguration()
	game := gameplay.NewGame(config.GuessLimit, gameplay.NewGameSettings(config.GuessLimit, config.MinNum, config.MaxNum))
	guessBot := tgbot.NewBot(updates, bot, game)

	guessBot.Run(ctx)
}
