package gameplay

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type GameSettings struct {
	guessLimit int
	minNum     int
	maxNum     int
}

type Game struct {
	guessLeft int
	target    int
	settings  *GameSettings
}

func NewGame(guess int, set *GameSettings) *Game {
	return &Game{
		guessLeft: guess,
		settings:  set,
	}
}

func NewGameSettings(limit int, min int, max int) *GameSettings {
	return &GameSettings{
		guessLimit: limit,
		minNum:     min,
		maxNum:     max,
	}
}

func (g *Game) HandleProcessMessage(userMessage string) (string, bool) {
	botAnswerOptions := map[string]string{
		"guessRange":     fmt.Sprintf("Я загадал число в диапазоне от %d до %d", g.settings.minNum, g.settings.maxNum),
		"notInt":         "Введено не число. Повтори попытку",
		"numberTooBig":   fmt.Sprintf("Твое число слишком БОЛЬШОЕ. Число оставшихся попыток: %d", g.guessLeft-1),
		"numberTooSmall": fmt.Sprintf("Твое число слишком МАЛЕНЬКОЕ. Число оставшихся попыток: %d", g.guessLeft-1),
		"playerWon":      "О, да! У тебя ПОЛУЧИЛОСЬ отгадать число!\nЧтобы сыграть еще жми /start",
		"playerLost":     fmt.Sprintf("Извини, у тебя не получилось отгадать число. Ответ был %d.\nЧтобы сыграть еще жми /start", g.target), //nolint:lll
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
	case userMessageAsInt > g.target:
		botAnswer = botAnswerOptions["numberTooBig"]

		g.guessLeft--
	case userMessageAsInt < g.target:
		botAnswer = botAnswerOptions["numberTooSmall"]

		g.guessLeft--
	case userMessageAsInt == g.target:
		botAnswer = botAnswerOptions["playerWon"]
	}

	if g.guessLeft == 0 {
		botAnswer = botAnswerOptions["playerLost"]

		g.guessLeft = g.settings.guessLimit
	}

	return botAnswer, showResponseKeyboard
}

func (g *Game) HandleProcessCallbackQuery(callbackData string) string {
	botAnswerOptions := map[string]string{
		"start":     fmt.Sprintf("Тогда погнали! У тебя %d попыток. Отправь мне число", g.settings.guessLimit),
		"rejection": "Если передумаешь, жми /start",
	}

	var botAnswer string

	switch {
	case callbackData == "yes":
		g.target = g.generateRandNumb()
		g.guessLeft = g.settings.guessLimit

		botAnswer = botAnswerOptions["start"]
	case callbackData == "no":
		botAnswer = botAnswerOptions["rejection"]
	}

	return botAnswer
}

func (g *Game) generateRandNumb() int {
	max := big.NewInt(int64(g.settings.maxNum))

	randomNumber, err := rand.Int(rand.Reader, max)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "gameplay",
			"function": "generateRandNumb",
			"error":    err,
		}).Warning("Random number generation error")
	}

	result := randomNumber.Int64() + int64(g.settings.minNum)

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
