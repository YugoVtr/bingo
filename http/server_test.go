package http_test

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yugovtr/bingo/http"
)

func TestNewServer(t *testing.T) {
	client := AssertServer(t)
	err := client.HealthCheck()
	assert.NoError(t, err)
}

func TestServer_WebSocket(t *testing.T) {
	client := AssertServer(t)
	err := client.WebSocket()
	assert.NoError(t, err)
}

func AssertServer(t *testing.T) *http.Client {
	t.Helper()

	server := http.NewServer(http.ServerConfig{})
	s := httptest.NewServer(server.Handler)

	t.Cleanup(func() { server.Close() })

	return http.NewClient(s.URL, t)
}
