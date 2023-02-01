package repository

import (
	"encoding/json"

	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const (
	DB = "bingo"
)

type Number struct {
	ID    string `gorethink:"id,omitempty"`
	Value int    `gorethink:"value"`
}

func NewRethinkDB(s *r.Session) *rethinkDB {
	return &rethinkDB{s}
}

type rethinkDB struct {
	session *r.Session
}

func (s *rethinkDB) AddHistoric(h int) {
	if err := s.db().Table("historic").Insert(Number{Value: h}).Exec(s.session); err != nil {
		panic(err)
	}
}

func (s *rethinkDB) ListenHistoric(ch chan<- int) {
	rows, err := s.db().Table("historic").Changes().Run(s.session)
	if err != nil {
		panic(err)
	}

	go func() {
		defer rows.Close()

		var change r.ChangeResponse
		for rows.Next(&change) {
			if newValue, ok := change.NewValue.(map[string]interface{}); ok {
				var number Number
				data, err := json.Marshal(newValue)
				if err != nil {
					continue
				}

				if err := json.Unmarshal(data, &number); err != nil {
					continue
				}

				ch <- number.Value
			}
		}
	}()
}

func (s *rethinkDB) db() r.Term {
	return r.DB(DB)
}
