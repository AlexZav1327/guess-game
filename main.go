package main

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var BOT_TOKEN = "5148810257:AAHyX2zm6Y154ioGsOPOp171sOWH7ZA924"

var target int
var guess int = 10

func main() {
	bot, err := tgbotapi.NewBotAPI(BOT_TOKEN)

	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message != nil {
			userMessage := update.Message.Text
			var response tgbotapi.Chattable = tgbotapi.NewMessage(update.Message.Chat.ID, "")
			var responseKeyboard tgbotapi.Chattable = tgbotapi.NewMessage(update.Message.Chat.ID, "")

			if userMessage == "/start" {
				response = tgbotapi.NewMessage(update.Message.Chat.ID, "Я загадал число в диапазоне от 1 до 100")
				responseKeyboard = createStartKeyboard(update.Message.Chat.ID)
			} else if convertUserMessageToInt(userMessage) > target {
				guess = guess - 1
				response = tgbotapi.NewMessage(update.Message.Chat.ID, "Твое число слишком БОЛЬШОЕ. Число оставшихся попыток: "+strconv.Itoa(guess))
			} else if convertUserMessageToInt(userMessage) < target {
				guess = guess - 1
				response = tgbotapi.NewMessage(update.Message.Chat.ID, "Твое число слишком МАЛЕНЬКОЕ. Число оставшихся попыток: "+strconv.Itoa(guess))
			} else if convertUserMessageToInt(userMessage) == target {
				response = tgbotapi.NewMessage(update.Message.Chat.ID, "О, да! У тебя ПОЛУЧИЛОСЬ отгадать число!\nЧтобы сыграть еще жми /start")
				seconds := time.Now().Unix()
				rand.Seed(seconds)
				target = rand.Intn(100) + 1
				guess = 10
			}

			if guess == 0 {
				response = tgbotapi.NewMessage(update.Message.Chat.ID, "Извини, у тебя не получилось отгадать число. Ответ был "+strconv.Itoa(target)+"."+"\nЧтобы сыграть еще жми /start")
				seconds := time.Now().Unix()
				rand.Seed(seconds)
				target = rand.Intn(100) + 1
				guess = 10
			}

			bot.Send(response)
			bot.Send(responseKeyboard)

		} else if update.CallbackQuery != nil {
			callbackData := update.CallbackQuery.Data
			var response tgbotapi.Chattable = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")

			if callbackData == "yes" {
				seconds := time.Now().Unix()
				rand.Seed(seconds)
				target = rand.Intn(100) + 1

				response = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Тогда погнали! Отправь мне число")

			} else if callbackData == "no" {
				response = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Если передумаешь, жми /start")
			}

			bot.Send(response)
		}
	}
}

func createStartKeyboard(chatId int64) tgbotapi.Chattable {
	btnY := tgbotapi.NewInlineKeyboardButtonData("Играем!", "yes")
	btnN := tgbotapi.NewInlineKeyboardButtonData("В другой раз", "no")

	row := tgbotapi.NewInlineKeyboardRow(btnY, btnN)
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(row)

	msg := tgbotapi.NewMessage(chatId, "Сможешь угадать это число с 10 попыток?")
	msg.ReplyMarkup = inlineKeyboard
	return msg
}

func convertUserMessageToInt(message string) int {
	res, err := strconv.Atoi(message)
	if err != nil {
		log.Fatal(err)
	}
	return res
}
