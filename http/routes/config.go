package routes

import (
	"net/http"

	"github.com/yugovtr/bingo/http/controllers"
	"github.com/yugovtr/bingo/http/handlers"
)

type Handler http.Handler

func New() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/healthcheck", handlers.HealthCheck)
	return router
}

func NewBingo(router *http.ServeMux, i controllers.NewBingoInput) *http.ServeMux {
	bingo := controllers.NewBingo(i)
	router.HandleFunc("/bingo/play", bingo.Play)
	router.HandleFunc("/bingo/next", bingo.Next)

	return router
}
