package game

import (
	"crypto/rand"
	"math"
	"math/big"
	"sort"

	"github.com/yugovtr/bingo/domain/entity"
)

const (
	MaxCardNumber = 75
	CardSize      = 25
)

type Bingo struct {
	generator func(_, _ int) int
	started   bool
}

func NewGame() *Bingo {
	b := &Bingo{generator: raffle}
	return b
}

// HasStarted checks if there has already been a number drawn for this game.
func (g *Bingo) HasStarted() bool {
	return g.started
}

// NewCard return a with card with pattern:
//
//	0:4 values between 1-15
//	5:9 values between 16-30
//	10:14 values between 31-45
//	15:19 values between 46-60
//	20:24 values between 61-75
func (g *Bingo) NewCard() (entity.Card, error) {
	if g.HasStarted() {
		return entity.Card{}, errGameStarted
	}

	drawnNumbers := map[int]struct{}{}
	card, lenght := make(entity.Card, CardSize), int(math.Sqrt(CardSize))

	for column := 0; column < lenght; column++ {
		for line := 0; line < lenght; line++ {

			start, end := (column*15)+1, ((column+1)*15)+1
			n := g.generator(start, end)

			for _, ok := drawnNumbers[n]; ok || (n < start || n >= end); _, ok = drawnNumbers[n] {
				// ignore repetitions
				n = g.generator(start, end)
			}

			drawnNumbers[n] = struct{}{}
			card[column*lenght+line] = n
		}
	}

	sort.Ints(card)
	return card, nil
}

// Raffle generate a number between 1 and MaxCardNumber.
//
//	Repetitions may occur in successive draws.
func (g *Bingo) Raffle() int {
	g.started = true

	return g.generator(0, MaxCardNumber)
}

func raffle(min, max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	return int(n.Int64()) + min
}
