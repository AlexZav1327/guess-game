package gameplay

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type GameSettings struct {
	GuessTotal int
	GuessLeft  int
	Target     int
	MinNum     int
	MaxNum     int
}

var StartingGameSettings = GameSettings{GuessTotal: 10, GuessLeft: 10, MinNum: 1, MaxNum: 1000}

func (g *GameSettings) HandleProcessMessage(userMessage string) (string, bool) {
	botAnswerOptions := map[string]string{
		"guessRange":     fmt.Sprintf("Я загадал число в диапазоне от %d до %d", g.MinNum, g.MaxNum),
		"notInt":         "Введено не число. Повтори попытку",
		"numberTooBig":   fmt.Sprintf("Твое число слишком БОЛЬШОЕ. Число оставшихся попыток: %d", g.GuessLeft-1),
		"numberTooSmall": fmt.Sprintf("Твое число слишком МАЛЕНЬКОЕ. Число оставшихся попыток: %d", g.GuessLeft-1),
		"playerWon":      "О, да! У тебя ПОЛУЧИЛОСЬ отгадать число!\nЧтобы сыграть еще жми /start",
		"playerLost":     fmt.Sprintf("Извини, у тебя не получилось отгадать число. Ответ был %d.\nЧтобы сыграть еще жми /start", g.Target), //nolint:lll
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
	case userMessageAsInt > g.Target:
		botAnswer = botAnswerOptions["numberTooBig"]

		g.GuessLeft--
	case userMessageAsInt < g.Target:
		botAnswer = botAnswerOptions["numberTooSmall"]

		g.GuessLeft--
	case userMessageAsInt == g.Target:
		botAnswer = botAnswerOptions["playerWon"]
	}

	if g.GuessLeft == 0 {
		botAnswer = botAnswerOptions["playerLost"]

		g.GuessLeft = g.GuessTotal
	}

	return botAnswer, showResponseKeyboard
}

func (g *GameSettings) HandleProcessCallbackQuery(callbackData string) string {
	botAnswerOptions := map[string]string{
		"start":     "Тогда погнали! Отправь мне число",
		"rejection": "Если передумаешь, жми /start",
	}

	var botAnswer string

	switch {
	case callbackData == "yes":
		g.Target = g.generateRandNumb()
		g.GuessLeft = g.GuessTotal

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
