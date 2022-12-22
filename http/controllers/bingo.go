package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/yugovtr/bingo/domain"
)

type Logger interface {
	Printf(string, ...any)
	Print(...any)
}

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Bingo struct {
	game        *domain.Game
	connections map[int]*websocket.Conn
	logger      Logger
}

func NewBingo(game *domain.Game) *Bingo {
	return &Bingo{
		game:        game,
		connections: make(map[int]*websocket.Conn),
		logger:      log.Default(),
	}
}

func (b *Bingo) Next(w http.ResponseWriter, r *http.Request) {
	n := b.game.Raffle()
	b.logger.Printf("new number requested: %d", n)

	w.Write(writeInt(n))

	if winner, ok := b.game.HasWinner(); ok {
		b.logger.Printf("we have a winner: %d", int(*winner))
		conn, ok := b.connections[int(*winner)]
		if !ok {
			return
		}

		conn.WriteMessage(1, writeString("you win"))
	}
}

func (b *Bingo) Play(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		b.logger.Printf("connection error: %s", err)
		return
	}

	myNumber := b.game.Play()
	b.connections[myNumber] = conn

	b.logger.Printf("new connection started with number %d", myNumber)
	conn.WriteMessage(1, writeInt(myNumber))
}

func writeInt(i int) []byte {
	return []byte(fmt.Sprintf("%d", i))
}

func writeString(s string) []byte {
	return []byte(fmt.Sprintf("%s", s))
}
