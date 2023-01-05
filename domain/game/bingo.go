package game

import (
	"crypto/rand"
	"math/big"

	"github.com/yugovtr/bingo/domain/entity"
	"golang.org/x/exp/slices"
)

const MaxCardNumber = 10

type Repository interface {
	AddHistoric(int)
	Historic() []int
	AddPlayer(entity.Player)
	Players() []entity.Player
}

type (
	Caller     func() int
	BingoError string
)

type Bingo struct {
	Repository
	Caller
}

func (err BingoError) Error() string {
	return string(err)
}

func Error(err string) error {
	return BingoError(err)
}

func NewGame(repo Repository) *Bingo {
	return &Bingo{Caller: raffle, Repository: repo}
}

func NewGameWithCaller(repo Repository, f func() int) *Bingo {
	return &Bingo{Caller: f, Repository: repo}
}

func (g *Bingo) HasWinner() (*entity.Player, bool) {
	var p entity.Player
	i := slices.IndexFunc(g.Players(), func(i entity.Player) bool {
		return slices.Index(g.Historic(), int(i)) >= 0
	})

	if i < 0 {
		return &p, false
	}

	p = entity.Player(g.Players()[i])
	return &p, i >= 0
}

func (g *Bingo) HasStarted() bool {
	return len(g.Historic()) > 0
}

func (g *Bingo) Play() (int, error) {
	if g.HasStarted() {
		return -1, Error("game already started")
	}

	var n int
	for n = g.Caller(); slices.Index(g.Players(), entity.Player(n)) != -1; n = g.Caller() {
		// non repet players
		// TODO - check if has a number to create a player
	}

	g.AddPlayer(entity.Player(n))
	return n, nil
}

func (g *Bingo) Raffle() (int, error) {
	if _, ok := g.HasWinner(); ok {
		return -1, Error("game already done")
	}

	var n int
	for n = g.Caller(); slices.Index(g.Historic(), n) != -1; n = g.Caller() {
		// non repet numbers
	}

	g.AddHistoric(n)
	return n, nil
}

func raffle() int {
	// [0, max) means greater than or equal to 0 and less than (max+1)
	n, _ := rand.Int(rand.Reader, big.NewInt(MaxCardNumber+1))
	return int(n.Int64())
}
