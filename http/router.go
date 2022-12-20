package http

import (
	"net/http"

	"github.com/yugovtr/bingo/http/handlers"
)

func newRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/", handlers.HealthCheck)
	router.HandleFunc("/healthcheck", handlers.HealthCheck)
	router.HandleFunc("/ws", handlers.WebSocket)

	return router
}
