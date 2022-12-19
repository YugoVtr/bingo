package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "embed"

	"github.com/gorilla/websocket"
	"github.com/yugovtr/bingo/internal/db"
)

const (
	dbName    = "chat"
	tableName = "messages"
)

//go:embed client.html
var server string

var (
	homeTemplate = template.Must(template.New("").Parse(server))
	addr         = flag.String("addr", "0.0.0.0:8081", "http service address")
	upgrader     = websocket.Upgrader{} // use default options
	dbClient     = new(db.DB)
	err          error
	pool         = []*websocket.Conn{}
)

type Model struct {
	User    string `json:"user"`
	Message string `json:"message"`
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	pool = append(pool, c)

	for {
		v := Model{}
		err := c.ReadJSON(&v)
		if err != nil {
			log.Println("read:", err)
			break
		}

		err = dbClient.Insert(tableName, map[string]string{
			"user":    v.User,
			"message": v.Message,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	flag.Parse()

	dbClient, err = db.Connect(dbName)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/echo", echo)
	http.HandleFunc("/truncate", func(w http.ResponseWriter, r *http.Request) {
		dbClient.Truncate(tableName)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		homeTemplate.Execute(w, fmt.Sprintf("ws://%s/echo", r.Host))
	})

	go dbClient.Listen(tableName, &pool)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
