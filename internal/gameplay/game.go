package gameplay

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type Game struct {
	GuessLeft int
	Target    int
}

type GameSettings struct {
	GuessLimit int
	MinNum     int
	MaxNum     int
}

func NewGame(guess int) *Game {
	return &Game{
		GuessLeft: guess,
	}
}

func NewGameSettings(limit int, min int, max int) *GameSettings {
	return &GameSettings{
		GuessLimit: limit,
		MinNum:     min,
		MaxNum:     max,
	}
}

func (g *GameSettings) HandleProcessMessage(userMessage string, target *int, guessLeft *int, guessLimit *int) (string, bool) { //nolint:lll
	botAnswerOptions := map[string]string{
		"guessRange":     fmt.Sprintf("Я загадал число в диапазоне от %d до %d", g.MinNum, g.MaxNum),
		"notInt":         "Введено не число. Повтори попытку",
		"numberTooBig":   fmt.Sprintf("Твое число слишком БОЛЬШОЕ. Число оставшихся попыток: %d", *guessLeft-1),
		"numberTooSmall": fmt.Sprintf("Твое число слишком МАЛЕНЬКОЕ. Число оставшихся попыток: %d", *guessLeft-1),
		"playerWon":      "О, да! У тебя ПОЛУЧИЛОСЬ отгадать число!\nЧтобы сыграть еще жми /start",
		"playerLost":     fmt.Sprintf("Извини, у тебя не получилось отгадать число. Ответ был %d.\nЧтобы сыграть еще жми /start", *target), //nolint:lll
	}

	var botAnswer string

	showResponseKeyboard := false
	userMessageAsInt, userMessageNotInt := convertUserMessageToInt(userMessage)

	switch {
	case userMessage == "/start":
		botAnswer = botAnswerOptions["guessRange"]
		showResponseKeyboard = true

	case userMessageNotInt:
		botAnswer = botAnswerOptions["notInt"]
	case userMessageAsInt > *target:
		botAnswer = botAnswerOptions["numberTooBig"]

		*guessLeft--
	case userMessageAsInt < *target:
		botAnswer = botAnswerOptions["numberTooSmall"]

		*guessLeft--
	case userMessageAsInt == *target:
		botAnswer = botAnswerOptions["playerWon"]
	}

	if *guessLeft == 0 {
		botAnswer = botAnswerOptions["playerLost"]

		*guessLeft = *guessLimit
	}

	return botAnswer, showResponseKeyboard
}

func (g *GameSettings) HandleProcessCallbackQuery(callbackData string, target *int, guessLeft *int, guessLimit *int) string { //nolint:lll
	botAnswerOptions := map[string]string{
		"start":     "Тогда погнали! Отправь мне число",
		"rejection": "Если передумаешь, жми /start",
	}

	var botAnswer string

	switch {
	case callbackData == "yes":
		*target = g.generateRandNumb()
		*guessLeft = *guessLimit

		botAnswer = botAnswerOptions["start"]
	case callbackData == "no":
		botAnswer = botAnswerOptions["rejection"]
	}

	return botAnswer
}

func (g *GameSettings) generateRandNumb() int {
	max := big.NewInt(int64(g.MaxNum))

	randomNumber, err := rand.Int(rand.Reader, max)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "gameplay",
			"function": "generateRandNumb",
			"error":    err,
		}).Warning("Random number generation error")
	}

	result := randomNumber.Int64() + int64(g.MinNum)

	return int(result)
}

func convertUserMessageToInt(message string) (int, bool) {
	notNumber := false

	result, err := strconv.Atoi(message)
	if err != nil {
		notNumber = true
	}

	return result, notNumber
}
