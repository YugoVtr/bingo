package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yugovtr/bingo/http"
	"github.com/yugovtr/bingo/infra/db"
)

func main() {
	addr := flag.String("addr", ":8081", "http service address")
	dbName := flag.String("dbName", "chat", "database name")

	flag.Parse()

	connection, err := db.Connect(*dbName)
	if err != nil {
		log.Fatal(err)
	}

	broadcast, err := connection.Listen("messages")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for b := range broadcast {
			log.Printf("%s\n", b)
		}
	}()

	log.Printf("server running in %s\n", *addr)
	defer log.Printf("server closed\n")

	config := http.ServerConfig{TCPAddress: *addr}

	select {
	case err := <-DelegateListenAndServe(http.NewServer(config).ListenAndServe):
		log.Fatal(err)
	case signal := <-NewInterruptSignal():
		log.Printf("received %s signal\n", signal)
	}
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
