package tgbot

import (
	"context"

	gameplay "github.com/AlexZav1327/guess-game/internal/gameplay"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

type GuessBot struct {
	updates tgbotapi.UpdatesChannel
	bot     tgbotapi.BotAPI
	game    gameplay.Game
}

func NewBot(updates tgbotapi.UpdatesChannel, bot *tgbotapi.BotAPI, game *gameplay.Game) *GuessBot {
	return &GuessBot{
		updates: updates,
		bot:     *bot,
		game:    *game,
	}
}

func (b *GuessBot) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case update := <-b.updates:
			switch {
			case update.Message != nil:
				b.processMessage(update)
			case update.CallbackQuery != nil:
				b.processCallbackQuery(update)
			}
		}
	}
}

func (b *GuessBot) processMessage(update tgbotapi.Update) {
	userMessage := update.Message.Text
	botAnswer, showResponseKeyboard := b.game.HandleProcessMessage(userMessage)
	response := tgbotapi.NewMessage(update.Message.Chat.ID, botAnswer)
	responseKeyboard := b.createStartKeyboard(update.Message.Chat.ID)

	_, err := b.bot.Send(response)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "tgbot",
			"function": "processMessage",
			"error":    err,
		}).Warning("Response sending error")

		return
	}

	if showResponseKeyboard {
		_, err := b.bot.Send(responseKeyboard)
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

func (b *GuessBot) processCallbackQuery(update tgbotapi.Update) {
	callbackData := update.CallbackQuery.Data
	botAswer := b.game.HandleProcessCallbackQuery(callbackData)
	response := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, botAswer)

	_, err := b.bot.Send(response)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "tgbot",
			"function": "processCallbackQuery",
			"error":    err,
		}).Warning("Response sending error")

		return
	}
}

func (b *GuessBot) createStartKeyboard(chatID int64) tgbotapi.Chattable {
	btnY := tgbotapi.NewInlineKeyboardButtonData("Играем!", "yes")
	btnN := tgbotapi.NewInlineKeyboardButtonData("В другой раз", "no")

	row := tgbotapi.NewInlineKeyboardRow(btnY, btnN)
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(row)

	msg := tgbotapi.NewMessage(chatID, "Сможешь угадать это число?")

	msg.ReplyMarkup = inlineKeyboard

	return msg
}
