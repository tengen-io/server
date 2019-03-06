package providers

import (
	"github.com/stretchr/testify/assert"
	"github.com/tengen-io/server/models"
	"github.com/tengen-io/server/test"
	"testing"
)

func TestGameProvider_GetGamesByIds(t *testing.T) {
	db := test.MakeDb()
	p := NewGameProvider(db)

	res, err := p.GetGamesByIds([]string{"1", "2"})
	assert.NoError(t, err)

	assert.Len(t, res, 2)
	assert.Equal(t, models.Invitation, res[0].State)
	assert.Equal(t, models.InProgress, res[1].State)
}

func TestGameProvider_GetGamesByState(t *testing.T) {
	db := test.MakeDb()
	p := NewGameProvider(db)

	res, err := p.GetGamesByState([]models.GameState{models.Invitation})
	assert.NoError(t, err)
	assert.True(t, len(res) > 0)
	assert.Len(t, res, 1)

	res, err = p.GetGamesByState([]models.GameState{models.Invitation, models.InProgress})
	assert.NoError(t, err)
	assert.True(t, len(res) > 1)
}

func TestGameProvider_CreateInvitation(t *testing.T) {
	db := test.MakeDb()
	p := NewGameProvider(db)

	res, err := p.CreateInvitation(models.CreateGameInvitationInput{
		BoardSize: 19,
		Type: models.Standard,
	})

	assert.NoError(t, err)
	assert.Equal(t, models.Standard, res.Type)
	assert.Equal(t, models.Invitation, res.State)
	assert.Equal(t, 19, res.BoardSize)
}
