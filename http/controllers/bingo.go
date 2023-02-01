package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/yugovtr/bingo/domain/contract"
)

type Logger interface {
	Printf(string, ...any)
	Print(...any)
}

type Event struct {
	Description string `json:"description"`
	Value       any    `json:"value"`
}

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Bingo struct {
	contract.Bingo
	Repository contract.BingoRepository
	logger     Logger
}

type NewBingoInput struct {
	Game       contract.Bingo
	Repository contract.BingoRepository
}

func NewBingo(i NewBingoInput) *Bingo {
	return &Bingo{
		Bingo:      i.Game,
		Repository: i.Repository,
		logger:     log.Default(),
	}
}

func (b *Bingo) Next(w http.ResponseWriter, r *http.Request) {
	n := b.Raffle()
	b.Repository.AddHistoric(n)

	b.logger.Printf("new number requested: %d", n)
	_, _ = w.Write(writeInt(n))
}

func (b *Bingo) Play(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		b.logger.Printf("connection error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	myCard, err := b.NewCard()
	if err != nil {
		conn.Close()
		b.logger.Printf("game error: %s", err)
		return
	}

	b.logger.Printf("new connection started with card %d", myCard)
	event, _ := json.Marshal(Event{
		Description: "card",
		Value:       myCard,
	})
	_ = conn.WriteMessage(1, event)

	ch := make(chan int)
	b.Repository.ListenHistoric(ch)

	go func() {
		defer close(ch)

		b.logger.Print("listening drawn numbers...")
		defer b.logger.Print("listening drawn numbers...done")

		for c := range ch {
			b.logger.Printf("new number drawn: %d", c)

			_ = conn.WriteJSON(Event{
				Description: "drawn",
				Value:       c,
			})
		}
	}()

	for {
		event := &Event{}
		err := conn.ReadJSON(event)
		if err != nil {
			switch err.(type) {
			case *websocket.CloseError:
				b.logger.Printf("client connection closed")
			default:
				b.logger.Printf("read message error: %s, %T", err, err)
			}

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		b.logger.Printf("new client event: %s", event.Description)
	}
}

func writeInt(i int) []byte {
	return []byte(fmt.Sprintf("%d", i))
}
