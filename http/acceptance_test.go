package http_test

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yugovtr/bingo/domain/entity"
	"github.com/yugovtr/bingo/http"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yugovtr/bingo/http/controllers"
)

func AssertServer(t *testing.T, config http.ServerConfig) *http.Client {
	t.Helper()

	server := http.NewServer(config)
	s := httptest.NewServer(server.Handler)

	t.Cleanup(func() { server.Close() })
	return http.NewClient(s.URL, t)
}

func ReadEvent(t *testing.T, p *websocket.Conn) (*controllers.Event, error) {
	t.Helper()

	event := make(chan *controllers.Event)
	go func() {
		e := &controllers.Event{}
		_ = p.ReadJSON(e)
		event <- e
	}()

	select {
	case e := <-event:
		return e, nil
	case <-time.After(2 * time.Second):
		return &controllers.Event{}, errors.New("socket read timeout")
	}
}

func AcceptanceBingo(t *testing.T, client *http.Client) {
	t.Helper()

	// requests to participate in the game
	player, err := client.PlayBingo()
	assert.NoError(t, err)
	defer player.Close()

	// wait for the card
	event, err := ReadEvent(t, player)
	require.NoError(t, err)
	assert.Equal(t, "card", event.Description)

	// asks for a new number
	b := bytes.Buffer{}
	err = client.BingoNext(&b)
	assert.NoError(t, err)
	assert.Equal(t, "1", b.String())

	_ = client.BingoNext(&bytes.Buffer{}) // TODO: if call only one time break this test

	// wait to receive the number drawn
	event, err = ReadEvent(t, player)
	require.NoError(t, err)
	assert.Equal(t, "drawn", event.Description)

	// 	bingo!!!
	err = player.WriteJSON(controllers.Event{Description: "bingo", Value: ""})
	assert.NoError(t, err)
}

/**** STUB ****/
type StubGame struct{}

func (s StubGame) HasStarted() bool {
	return false
}

func (s StubGame) NewCard() (entity.Card, error) {
	return make(entity.Card, 25), nil
}

func (s StubGame) Raffle() int {
	return 1
}
