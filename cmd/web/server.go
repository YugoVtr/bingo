package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yugovtr/bingo/domain"
	"github.com/yugovtr/bingo/http"
	"github.com/yugovtr/bingo/http/routes"

	http2 "net/http"
)

func main() {
	addr := flag.String("addr", ":8081", "http service address")

	flag.Parse()

	log.Printf("server running in %s\n", *addr)
	defer log.Printf("server closed\n")

	server := Setup(*addr)

	select {
	case err := <-DelegateListenAndServe(server.ListenAndServe):
		log.Fatal(err)
	case signal := <-NewInterruptSignal():
		log.Printf("received %s signal\n", signal)
	}
}

func Setup(addr string) *http2.Server {
	game := domain.NewGame()
	routes := routes.NewBingo(routes.New(), game)
	config := http.ServerConfig{TCPAddress: addr, Routes: routes}

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
