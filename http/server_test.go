package http_test

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yugovtr/bingo/domain/repository"
	"github.com/yugovtr/bingo/http"
	"github.com/yugovtr/bingo/http/controllers"
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
	log.SetPrefix("[SERVER] ")

	input := controllers.NewBingoInput{
		Game:       &StubGame{},
		Repository: repository.NewInMemory(),
	}

	routes := routes.NewBingo(routes.New(), input)

	client := AssertServer(t, http.ServerConfig{Routes: routes})
	AcceptanceBingo(t, client)
}
