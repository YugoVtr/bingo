package domain

import (
	"crypto/rand"
	"math/big"

	"golang.org/x/exp/slices"
)

const MaxCardNumber = 10

type Caller func() int

type Player int

type BingoError string

func (err BingoError) Error() string {
	return string(err)
}

func Error(err string) error {
	return BingoError(err)
}

type Game struct {
	Caller
	History []int
	Players []Player
}

func NewGame() *Game {
	return &Game{Caller: raffle}
}

func NewGameWithCaller(f func() int) *Game {
	return &Game{Caller: f}
}

func (g *Game) HasWinner() (*Player, bool) {
	var p Player
	i := slices.IndexFunc(g.Players, func(i Player) bool {
		return slices.Index(g.History, int(i)) >= 0
	})

	if i < 0 {
		return &p, false
	}

	p = Player(g.Players[i])
	return &p, i >= 0
}

func (g *Game) HasStarted() bool {
	return len(g.History) > 0
}

func (g *Game) Play() (int, error) {
	if g.HasStarted() {
		return -1, Error("game already started")
	}

	var n int
	for n = g.Caller(); slices.Index(g.Players, Player(n)) != -1; n = g.Caller() {
		// non repet players
		// TODO - check if has a number to create a player
	}

	g.Players = append(g.Players, Player(n))
	return n, nil
}

func (g *Game) Raffle() (int, error) {
	if _, ok := g.HasWinner(); ok {
		return -1, Error("game already done")
	}

	var n int
	for n = g.Caller(); slices.Index(g.History, n) != -1; n = g.Caller() {
		// non repet numbers
	}

	g.History = append(g.History, n)
	return n, nil
}

func raffle() int {
	// [0, max) means greater than or equal to 0 and less than (max+1)
	n, _ := rand.Int(rand.Reader, big.NewInt(MaxCardNumber+1))
	return int(n.Int64())
}
