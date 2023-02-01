package contract

import "github.com/yugovtr/bingo/domain/entity"

type Bingo interface {
	HasStarted() bool
	NewCard() (entity.Card, error)
	Raffle() int
}
