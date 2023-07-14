package gameplay

import (
	"crypto/rand"
	"math/big"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type GameSettings struct {
	GuessQty int
	Target   int
	MinNum   int
	MaxNum   int
}

func (g *GameSettings) GenerateRandNumb() int {
	max := big.NewInt(int64(g.MaxNum))

	randomNumber, err := rand.Int(rand.Reader, max)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "gameplay",
			"function": "GenerateRandNumb",
			"error":    err,
		}).Warning("Random number generation error")
	}

	result := randomNumber.Int64() + int64(g.MinNum)

	return int(result)
}

func (GameSettings) ConvertUserMessageToInt(message string) (int, bool) {
	notNumber := false

	result, err := strconv.Atoi(message)
	if err != nil {
		notNumber = true
	}

	return result, notNumber
}
