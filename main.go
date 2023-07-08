package main

import (
	"context"
	"fmt"
	"math/rand"
	"os/signal"
	"strconv"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	const botToken = "5148810257:AAHyX2zm6Y154ioGsOPOp171sOWH7ZA924"

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		fmt.Println("BotAPI instance creation error:", err)
		return
	}

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := bot.GetUpdatesChan(updateConfig)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	defer cancel()

	run(ctx, bot, updates)
}

func run(ctx context.Context, bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) {
	var target int
	guess := 10

	for {
		select {
		case <-ctx.Done():
			return
		case update := <-updates:
			switch {
			case update.Message != nil:
				processMessage(&target, &guess, bot, update)
			case update.CallbackQuery != nil:
				processCallbackQuery(&target, bot, update)
			}
		}
	}
}

func processMessage(target *int, guess *int, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userMessage := update.Message.Text
	var response tgbotapi.Chattable = tgbotapi.NewMessage(update.Message.Chat.ID, "")
	var responseKeyboard tgbotapi.Chattable = tgbotapi.NewMessage(update.Message.Chat.ID, "")

	switch {
	case userMessage == "/start":
		response = tgbotapi.NewMessage(update.Message.Chat.ID, "Я загадал число в диапазоне от 1 до 100")
		responseKeyboard = createStartKeyboard(update.Message.Chat.ID)
	case convertUserMessageToInt(userMessage) > *target:
		*guess = *guess - 1
		response = tgbotapi.NewMessage(update.Message.Chat.ID, "Твое число слишком БОЛЬШОЕ. Число оставшихся попыток: "+strconv.Itoa(*guess))
	case convertUserMessageToInt(userMessage) < *target:
		*guess = *guess - 1
		response = tgbotapi.NewMessage(update.Message.Chat.ID, "Твое число слишком МАЛЕНЬКОЕ. Число оставшихся попыток: "+strconv.Itoa(*guess))
	case convertUserMessageToInt(userMessage) == *target:
		response = tgbotapi.NewMessage(update.Message.Chat.ID, "О, да! У тебя ПОЛУЧИЛОСЬ отгадать число!\nЧтобы сыграть еще жми /start")
		*target = rand.Intn(100) + 1
		*guess = 10
	}

	if *guess == 0 {
		response = tgbotapi.NewMessage(update.Message.Chat.ID, "Извини, у тебя не получилось отгадать число. Ответ был "+strconv.Itoa(*target)+"."+"\nЧтобы сыграть еще жми /start")
		*target = rand.Intn(100) + 1
		*guess = 10
	}

	bot.Send(response)
	bot.Send(responseKeyboard)
}

func processCallbackQuery(target *int, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	callbackData := update.CallbackQuery.Data
	var response tgbotapi.Chattable = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")

	switch {
	case callbackData == "yes":
		*target = rand.Intn(100) + 1
		response = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Тогда погнали! Отправь мне число")
	case callbackData == "no":
		response = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Если передумаешь, жми /start")
	}

	bot.Send(response)
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
	result, err := strconv.Atoi(message)
	if err != nil {
		fmt.Println("Error converting user response to integer:", err)
		return 0
	}

	return result
}
