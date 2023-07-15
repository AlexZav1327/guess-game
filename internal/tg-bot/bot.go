package tgbot

import (
	"context"
	"fmt"

	gameplay "github.com/AlexZav1327/guess-game/internal/gameplay"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

type GuessBot struct {
	Updates tgbotapi.UpdatesChannel
	Bot     tgbotapi.BotAPI
}

func (b *GuessBot) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case update := <-b.Updates:
			switch {
			case update.Message != nil:
				b.processMessage(update, &gameplay.StartingGameSettings)
			case update.CallbackQuery != nil:
				b.processCallbackQuery(update, &gameplay.StartingGameSettings)
			}
		}
	}
}

func (b *GuessBot) processMessage(update tgbotapi.Update, game *gameplay.GameSettings) {
	userMessage := update.Message.Text
	botAnswer, showResponseKeyboard := game.HandleProcessMessage(userMessage)
	response := tgbotapi.NewMessage(update.Message.Chat.ID, botAnswer)
	responseKeyboard := createStartKeyboard(update.Message.Chat.ID, game.GuessTotal)

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

func (b *GuessBot) processCallbackQuery(update tgbotapi.Update, game *gameplay.GameSettings) {
	callbackData := update.CallbackQuery.Data
	botAswer := game.HandleProcessCallbackQuery(callbackData)
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
