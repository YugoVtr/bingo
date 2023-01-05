package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/yugovtr/bingo/domain/game"
)

type Logger interface {
	Printf(string, ...any)
	Print(...any)
}

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Bingo struct {
	game        *game.Bingo
	connections map[int]*websocket.Conn
	logger      Logger
}

func NewBingo(game *game.Bingo) *Bingo {
	return &Bingo{
		game:        game,
		connections: make(map[int]*websocket.Conn),
		logger:      log.Default(),
	}
}

func (b *Bingo) Next(w http.ResponseWriter, r *http.Request) {
	n, err := b.game.Raffle()
	if err != nil {
		b.logger.Printf("raffle error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b.logger.Printf("new number requested: %d", n)
	_, _ = w.Write(writeInt(n))

	if winner, ok := b.game.HasWinner(); ok {
		b.logger.Printf("we have a winner: %d", int(*winner))

		conn, ok := b.connections[int(*winner)]
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_ = conn.WriteMessage(1, writeString("you win"))
	}
}

func (b *Bingo) Play(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		b.logger.Printf("connection error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	myNumber, err := b.game.Play()
	if err != nil {
		conn.Close()
		b.logger.Printf("game error: %s", err)
		return
	}

	b.connections[myNumber] = conn

	b.logger.Printf("new connection started with number %d", myNumber)
	_ = conn.WriteMessage(1, writeInt(myNumber))
}

func writeInt(i int) []byte {
	return []byte(fmt.Sprintf("%d", i))
}

func writeString(s string) []byte {
	return []byte(s)
}
