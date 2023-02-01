package game_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	bingo "github.com/yugovtr/bingo/domain/game"
)

func TestBingo_NewCard(t *testing.T) {
	session := AssertNewGame(t)

	t.Run("when has new valid card", func(t *testing.T) {
		card, err := session.NewCard()
		assert.NoError(t, err)
		assert.Len(t, card, 25)
		t.Log(card)
	})

	t.Run("when game already started", func(t *testing.T) {
		_ = session.Raffle()
		card, err := session.NewCard()
		assert.Error(t, err)
		t.Log(card)
	})
}

func AssertNewGame(t *testing.T) *bingo.Bingo {
	t.Helper()

	session := bingo.NewGame()
	assert.NotNil(t, session)

	return session
}
