package main

import (
	"log"

	"github.com/yugovtr/bingo/client/db"
)

func main() {
	const (
		dbName    = "chat"
		tableName = "messages"
	)

	cnt, err := db.Connect(dbName)
	if err != nil {
		log.Fatal(err)
	}

	cnt.Listen(tableName)
}
