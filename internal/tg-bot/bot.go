package tgbot

import (
	"context"
	"fmt"

	gameplay "github.com/AlexZav1327/guess-game/internal/gameplay"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

type GuessBot struct {
	Updates  tgbotapi.UpdatesChannel
	Bot      tgbotapi.BotAPI
	Game     gameplay.Game
	Settings gameplay.GameSettings
}

func NewBot(updates tgbotapi.UpdatesChannel, bot tgbotapi.BotAPI, game gameplay.Game, settings gameplay.GameSettings) *GuessBot { //nolint:lll
	return &GuessBot{
		Updates:  updates,
		Bot:      bot,
		Game:     game,
		Settings: settings,
	}
}

func (b *GuessBot) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case update := <-b.Updates:
			switch {
			case update.Message != nil:
				b.processMessage(update, &b.Game.Target, &b.Game.GuessLeft, &b.Settings.GuessLimit)
			case update.CallbackQuery != nil:
				b.processCallbackQuery(update, &b.Game.Target, &b.Game.GuessLeft, &b.Settings.GuessLimit)
			}
		}
	}
}

func (b *GuessBot) processMessage(update tgbotapi.Update, target *int, guessLeft *int, guessLimit *int) {
	userMessage := update.Message.Text
	botAnswer, showResponseKeyboard := b.Settings.HandleProcessMessage(userMessage, target, guessLeft, guessLimit)
	response := tgbotapi.NewMessage(update.Message.Chat.ID, botAnswer)
	responseKeyboard := createStartKeyboard(update.Message.Chat.ID, *guessLimit)

	_, err := b.Bot.Send(response)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "tgbot",
			"function": "processMessage",
			"error":    err,
		}).Warning("Response sending error")

		return
	}

	if showResponseKeyboard {
		_, err := b.Bot.Send(responseKeyboard)
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "tgbot",
				"function": "processMessage",
				"error":    err,
			}).Warning("Response keyboard sending error")

			return
		}
	}
}

func (b *GuessBot) processCallbackQuery(update tgbotapi.Update, target *int, guessLeft *int, guessLimit *int) {
	callbackData := update.CallbackQuery.Data
	botAswer := b.Settings.HandleProcessCallbackQuery(callbackData, target, guessLeft, guessLimit)
	response := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, botAswer)

	_, err := b.Bot.Send(response)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "tgbot",
			"function": "processCallbackQuery",
			"error":    err,
		}).Warning("Response sending error")

		return
	}
}

func createStartKeyboard(chatID int64, guessTotal int) tgbotapi.Chattable {
	btnY := tgbotapi.NewInlineKeyboardButtonData("Играем!", "yes")
	btnN := tgbotapi.NewInlineKeyboardButtonData("В другой раз", "no")

	row := tgbotapi.NewInlineKeyboardRow(btnY, btnN)
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(row)

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Сможешь угадать это число с %d попыток?", guessTotal))
	msg.ReplyMarkup = inlineKeyboard

	return msg
}
