package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yugovtr/bingo/domain/game"
	"github.com/yugovtr/bingo/domain/repository"
	"github.com/yugovtr/bingo/http"
	"github.com/yugovtr/bingo/http/routes"
	database "github.com/yugovtr/bingo/infra/db"

	http2 "net/http"
)

func main() {
	host := flag.String("host", ":8081", "http service address")
	db := flag.String("db", ":28015", "rethinkdb address")

	flag.Parse()

	log.Printf("server running in %s\n", *host)
	defer log.Printf("server closed\n")

	server := Setup(*host, *db)

	select {
	case err := <-DelegateListenAndServe(server.ListenAndServe):
		log.Fatal(err)
	case signal := <-NewInterruptSignal():
		log.Printf("received %s signal\n", signal)
	}
}

func Setup(host, db string) *http2.Server {
	const timeout = 10 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client := database.Connect(ctx, db)
	game := game.NewGame(repository.NewRethinkDB(client.Session))
	routes := routes.NewBingo(routes.New(), game)
	config := http.ServerConfig{TCPAddress: host, Routes: routes}

	return http.NewServer(config)
}

func DelegateListenAndServe(serve func() error) chan error {
	listenErr := make(chan error)

	go func() {
		if err := serve(); err != nil {
			listenErr <- err
		}
	}()

	return listenErr
}

func NewInterruptSignal() <-chan os.Signal {
	signals := make(chan os.Signal, 1)
	signalsToListen := []os.Signal{
		syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM,
	}

	signal.Notify(signals, signalsToListen...)
	return signals
}
