package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/yugovtr/bingo/client/db"
)

/*
	r.dbCreate('chat');
	r.db('chat').tableCreate('messages')
*/

func prompt(in string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s ", in)
	out, _ := reader.ReadString('\n')
	return strings.TrimSuffix(out, "\n")
}

func main() {
	const (
		dbName    = "chat"
		tableName = "messages"
	)

	c, err := db.Connect(dbName)
	if err != nil {
		log.Fatal(err)
	}

	usr := prompt("usr:")

	for msg := prompt(">"); len(msg) > 0; msg = prompt(">") {
		err = c.Insert(tableName, map[string]string{
			"user":    usr,
			"message": msg,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}
