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
	myNumber := session.Play()

	newNumber := session.Raffle()
	for myNumber != newNumber {
		newNumber = session.Raffle()
	}

	winner, ok := session.HasWinner()

	assert.NotNil(t, winner)
	assert.True(t, ok, "winner not found")

	assert.Equal(t, myNumber, newNumber)
	t.Logf("win after %d numbers", len(session.History))
}

func AssertNewGame(t *testing.T) *domain.Game {
	t.Helper()

	session := domain.NewGame()
	assert.NotNil(t, session)

	return session
}
