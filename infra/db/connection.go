package db

import (
	"fmt"

	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

/*
	r.dbCreate('chat');
	r.db('chat').tableCreate('messages');
	r.db('chat').table("messages").insert({
    user: "vitor",
    message: "Dolor sit amet"
	})
*/

type DB struct {
	dbName  string
	session *r.Session
}

func Connect(dbName string) (db *DB, err error) {
	session, err := r.Connect(r.ConnectOpts{
		Address: "0.0.0.0:28015",
	})
	if err != nil {
		return nil, fmt.Errorf("connection error: %w", err)
	}

	return &DB{dbName: dbName, session: session}, nil
}

func (db *DB) Insert(table string, value any) error {
	err := r.DB(db.dbName).Table(table).Insert(value).Exec(db.session)
	return err
}

func (db *DB) Truncate(table string) error {
	err := r.DB(db.dbName).Table(table).Delete().Exec(db.session)
	return err
}

func (db *DB) Listen(table string) (<-chan []byte, error) {
	rows, err := r.DB(db.dbName).Table(table).Changes().Run(db.session)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	broadcast := make(chan []byte)

	ch := make(chan map[string]map[string]string)
	rows.Listen(ch)

	go func() {
		defer close(broadcast)

		for c := range ch {
			v := c["new_val"]
			if len(v["user"]) == 0 || len(v["message"]) == 0 {
				continue
			}

			msg := []byte(fmt.Sprintf(`{"user": "%s", "message": "%s"}`, v["user"], v["message"]))
			broadcast <- msg
		}
	}()

	return broadcast, nil
}
