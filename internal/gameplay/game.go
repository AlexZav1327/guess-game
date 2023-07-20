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

type DefaultConfiguration struct {
	GuessLimit int
	MinNum     int
	MaxNum     int
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

func NewDefaultConfiguration() *DefaultConfiguration {
	return &DefaultConfiguration{
		GuessLimit: 10,
		MinNum:     1,
		MaxNum:     1000,
	}
}

var botAnswerOptions = map[string]string{ //nolint:gochecknoglobals
	"guessRange":     "Я загадал число в диапазоне от %d до %d",
	"notInt":         "Введено не число. Повтори попытку",
	"numberTooBig":   "Твое число слишком БОЛЬШОЕ. Число оставшихся попыток: %d",
	"numberTooSmall": "Твое число слишком МАЛЕНЬКОЕ. Число оставшихся попыток: %d",
	"playerWon":      "О, да! У тебя ПОЛУЧИЛОСЬ отгадать число!\nЧтобы сыграть еще жми /start",
	"playerLost":     "Извини, у тебя не получилось отгадать число. Ответ был %d.\nЧтобы сыграть еще жми /start",
	"start":          "Тогда погнали! У тебя %d попыток. Отправь мне число",
	"rejection":      "Если передумаешь, жми /start",
}

func (g *Game) HandleProcessMessage(userMessage string) (string, bool) {
	var botAnswer string

	showResponseKeyboard := false
	userMessageAsInt, err := strconv.Atoi(userMessage)

	switch {
	case userMessage == "/start":
		botAnswer = fmt.Sprintf(botAnswerOptions["guessRange"], g.settings.minNum, g.settings.maxNum)
		showResponseKeyboard = true
	case err != nil:
		botAnswer = botAnswerOptions["notInt"]
	case userMessageAsInt > g.target:
		g.guessLeft--

		botAnswer = fmt.Sprintf(botAnswerOptions["numberTooBig"], g.guessLeft)
	case userMessageAsInt < g.target:
		g.guessLeft--

		botAnswer = fmt.Sprintf(botAnswerOptions["numberTooSmall"], g.guessLeft)
	case userMessageAsInt == g.target:
		botAnswer = botAnswerOptions["playerWon"]
	}

	if g.guessLeft == 0 {
		botAnswer = fmt.Sprintf(botAnswerOptions["playerLost"], g.target)
		g.guessLeft = g.settings.guessLimit
	}

	return botAnswer, showResponseKeyboard
}

func (g *Game) HandleProcessCallbackQuery(callbackData string) string {
	var botAnswer string

	switch {
	case callbackData == "yes":
		g.target = g.generateRandNum()
		g.guessLeft = g.settings.guessLimit

		botAnswer = fmt.Sprintf(botAnswerOptions["start"], g.settings.guessLimit)
	case callbackData == "no":
		botAnswer = botAnswerOptions["rejection"]
	}

	return botAnswer
}

func (g *Game) generateRandNum() int {
	max := big.NewInt(int64(g.settings.maxNum))

	randomNumber, err := rand.Int(rand.Reader, max)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "gameplay",
			"function": "generateRandNum",
			"error":    err,
		}).Panic("Random number generation error")
	}

	result := randomNumber.Int64() + int64(g.settings.minNum)

	return int(result)
}
