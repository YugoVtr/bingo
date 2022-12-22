package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yugovtr/bingo/domain"
)

func TestNewGame(t *testing.T) {
	AssertNewGame(t)
}

func TestGame_Play(t *testing.T) {
	session := AssertNewGame(t)
	myNumber, _ := session.Play()

	newNumber, _ := session.Raffle()
	for myNumber != newNumber {
		newNumber, _ = session.Raffle()
	}

	winner, ok := session.HasWinner()

	assert.NotNil(t, winner)
	assert.True(t, ok, "winner not found")

	assert.Equal(t, myNumber, newNumber)
	t.Logf("win after %d numbers", len(session.History))

	_, err := session.Play()
	assert.Error(t, err)

	_, err = session.Raffle()
	assert.Error(t, err)
}

func AssertNewGame(t *testing.T) *domain.Game {
	t.Helper()

	session := domain.NewGame()
	assert.NotNil(t, session)

	return session
}
