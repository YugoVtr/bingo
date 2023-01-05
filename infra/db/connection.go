package db

import (
	"context"
	"errors"
	"log"
	"time"

	"gopkg.in/rethinkdb/rethinkdb-go.v6"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type DB struct {
	Session *r.Session
}

func connect(addr string) <-chan *r.Session {
	var (
		ch = make(chan *r.Session)
	)

	go func() {
		for interval := time.Second; ; interval *= 2 {
			session, err := r.Connect(r.ConnectOpts{
				Address: addr,
			})
			if errors.As(err, &rethinkdb.RQLConnectionError{}) {
				log.Printf("rethinkdb connection refused in %s. retrying in %s...", addr, interval)
				time.Sleep(interval)
				continue
			}

			log.Printf("rethinkdb connected in %s", addr)
			ch <- session
			return
		}
	}()

	return ch
}

func Connect(ctx context.Context, address string) *DB {
	var s *r.Session

	select {
	case s = <-connect(address):
	case <-ctx.Done():
		log.Fatalf("connection attempt to database fail")
	}

	return &DB{Session: s}
}

/*
r.dbCreate('bingo');
r.db('bingo').tableCreate('historic');
r.db('bingo').tableCreate('player');

r.db('bingo').table("player").delete();
*/
func Migrate(s *r.Session) error {
	if err := r.DBCreate("bingo").Exec(s); err != nil {
		return err
	}

	if err := r.DB("bingo").TableCreate("historic").Exec(s); err != nil {
		return err
	}

	if err := r.DB("bingo").TableCreate("player").Exec(s); err != nil {
		return err
	}

	return nil
}
