package http

import (
	"net/http"

	"github.com/yugovtr/bingo/http/routes"
)

func NewServer(config ServerConfig) *http.Server {
	return &http.Server{
		Addr:    config.TCPAddress,
		Handler: config.Routes,
	}
}

type Logger interface {
	Log(...any)
	Logf(string, ...any)
}

type HTTPError error

type ServerConfig struct {
	TCPAddress string
	Routes     routes.Handler
}
