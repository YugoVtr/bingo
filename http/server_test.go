package http_test

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yugovtr/bingo/domain"
	"github.com/yugovtr/bingo/http"
	"github.com/yugovtr/bingo/http/routes"
)

func TestNewServer(t *testing.T) {
	client := AssertServer(t, http.ServerConfig{Routes: routes.New()})
	err := client.HealthCheck()
	assert.NoError(t, err)

	err = client.Home()
	assert.Error(t, err)
}

func TestServer_Bingo(t *testing.T) {
	game := domain.NewGameWithCaller(func() int { return 1 })
	routes := routes.NewBingo(routes.New(), game)

	client := AssertServer(t, http.ServerConfig{Routes: routes})

	player, err := client.PlayBingo()
	assert.NoError(t, err)
	defer player.Close()

	_, msg, err := player.ReadMessage()
	assert.NoError(t, err)
	assert.NotNil(t, msg)

	b := bytes.Buffer{}
	err = client.BingoNext(&b)

	assert.NoError(t, err)
	assert.Equal(t, "1", b.String())

	_, msg, err = player.ReadMessage()
	assert.NoError(t, err)
	assert.NotNil(t, msg)
	t.Logf("%s", msg)
}

func AssertServer(t *testing.T, config http.ServerConfig) *http.Client {
	t.Helper()

	server := http.NewServer(config)
	s := httptest.NewServer(server.Handler)

	t.Cleanup(func() { server.Close() })
	return http.NewClient(s.URL, t)
}
