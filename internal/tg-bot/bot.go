package tgbot

import (
	"context"
	"fmt"

	gameplay "github.com/AlexZav1327/guess-game/internal/gameplay"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

type Bot struct {
	Updates tgbotapi.UpdatesChannel
	Bot     tgbotapi.BotAPI
}

func (b *Bot) Run(ctx context.Context, game gameplay.GameSettings) {
	guessQty := game.GuessQty

	for {
		select {
		case <-ctx.Done():
			return
		case update := <-b.Updates:
			switch {
			case update.Message != nil:
				b.processMessage(&game.Target, &guessQty, game, update)
			case update.CallbackQuery != nil:
				b.processCallbackQuery(&game.Target, &guessQty, game, update)
			}
		}
	}
}

func (b *Bot) processMessage(target *int, guessQty *int, game gameplay.GameSettings, update tgbotapi.Update) {
	var response tgbotapi.Chattable = tgbotapi.NewMessage(update.Message.Chat.ID, "")

	var responseKeyboard tgbotapi.Chattable = tgbotapi.NewMessage(update.Message.Chat.ID, "")

	userMessage := update.Message.Text
	userMessageAsInt, userMessageNotInt := game.ConvertUserMessageToInt(userMessage)
	showResponseKeyboard := false

	botGameplayAnswers := map[string]string{
		"guessRange":     fmt.Sprintf("Я загадал число в диапазоне от %d до %d", game.MinNum, game.MaxNum),
		"notInt":         "Введено не число. Повтори попытку",
		"numberTooBig":   fmt.Sprintf("Твое число слишком БОЛЬШОЕ. Число оставшихся попыток: %d", *guessQty-1),
		"numberTooSmall": fmt.Sprintf("Твое число слишком МАЛЕНЬКОЕ. Число оставшихся попыток: %d", *guessQty-1),
		"playerWon":      "О, да! У тебя ПОЛУЧИЛОСЬ отгадать число!\nЧтобы сыграть еще жми /start",
		"playerLost":     fmt.Sprintf("Извини, у тебя не получилось отгадать число. Ответ был %d.\nЧтобы сыграть еще жми /start", *target), //nolint:lll
	}

	switch {
	case userMessage == "/start":
		response = tgbotapi.NewMessage(update.Message.Chat.ID, botGameplayAnswers["guessRange"])
		responseKeyboard = createStartKeyboard(update.Message.Chat.ID, game)
		showResponseKeyboard = true

	case userMessageNotInt:
		response = tgbotapi.NewMessage(update.Message.Chat.ID, botGameplayAnswers["notInt"])
	case userMessageAsInt > *target:
		*guessQty--

		response = tgbotapi.NewMessage(update.Message.Chat.ID, botGameplayAnswers["numberTooBig"])
	case userMessageAsInt < *target:
		*guessQty--

		response = tgbotapi.NewMessage(update.Message.Chat.ID, botGameplayAnswers["numberTooSmall"])
	case userMessageAsInt == *target:
		response = tgbotapi.NewMessage(update.Message.Chat.ID, botGameplayAnswers["playerWon"])
	}

	if *guessQty == 0 {
		response = tgbotapi.NewMessage(update.Message.Chat.ID, botGameplayAnswers["playerLost"])

		*guessQty = game.GuessQty
	}

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

func (b *Bot) processCallbackQuery(target *int, guessQty *int, game gameplay.GameSettings, update tgbotapi.Update) {
	var response tgbotapi.Chattable = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")

	callbackData := update.CallbackQuery.Data

	botGameplayAnswers := map[string]string{
		"start":     "Тогда погнали! Отправь мне число",
		"rejection": "Если передумаешь, жми /start",
	}

	switch {
	case callbackData == "yes":
		*target = game.GenerateRandNumb()
		*guessQty = game.GuessQty
		response = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, botGameplayAnswers["start"])
	case callbackData == "no":
		response = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, botGameplayAnswers["rejection"])
	}

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

func createStartKeyboard(chatID int64, game gameplay.GameSettings) tgbotapi.Chattable {
	botGameplayAnswers := map[string]string{
		"guessQty": fmt.Sprintf("Сможешь угадать это число с %d попыток?", game.GuessQty),
	}

	btnY := tgbotapi.NewInlineKeyboardButtonData("Играем!", "yes")
	btnN := tgbotapi.NewInlineKeyboardButtonData("В другой раз", "no")

	row := tgbotapi.NewInlineKeyboardRow(btnY, btnN)
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(row)

	msg := tgbotapi.NewMessage(chatID, botGameplayAnswers["guessQty"])
	msg.ReplyMarkup = inlineKeyboard

	return msg
}
