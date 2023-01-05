package repository

import (
	"github.com/yugovtr/bingo/domain/entity"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const (
	DB = "bingo"
)

type Number struct {
	Value int `json:"value"`
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

func (s *rethinkDB) Historic() []int {
	var result []int

	rows, err := s.db().Table("historic").Run(s.session)
	if err != nil {
		panic(err)
	}

	var numbers []Number
	if err := rows.All(&numbers); err != nil {
		panic(err)
	}

	for _, p := range numbers {
		result = append(result, p.Value)
	}

	return result
}

func (s *rethinkDB) AddPlayer(p entity.Player) {
	if err := s.db().Table("player").Insert(Number{Value: int(p)}).Exec(s.session); err != nil {
		panic(err)
	}
}

func (s *rethinkDB) Players() []entity.Player {
	var result []entity.Player

	rows, err := s.db().Table("player").Run(s.session)
	if err != nil {
		panic(err)
	}

	var numbers []Number
	if err := rows.All(&numbers); err != nil {
		panic(err)
	}

	for _, p := range numbers {
		result = append(result, entity.Player(p.Value))
	}

	return result
}

func (s *rethinkDB) db() r.Term {
	return r.DB(DB)
}
