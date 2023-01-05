package repository

import "github.com/yugovtr/bingo/domain/entity"

func NewInMemory() *gameInMemory {
	return new(gameInMemory)
}

type gameInMemory struct {
	historic []int
	players  []entity.Player
}

func (s *gameInMemory) AddHistoric(h int) {
	s.historic = append(s.historic, h)
}

func (s *gameInMemory) Historic() []int {
	return s.historic
}

func (s *gameInMemory) AddPlayer(p entity.Player) {
	s.players = append(s.players, p)
}

func (s *gameInMemory) Players() []entity.Player {
	return s.players
}
