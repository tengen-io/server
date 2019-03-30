package repository

import (
	"github.com/stretchr/testify/assert"
	"github.com/tengen-io/server/models"
	"github.com/tengen-io/server/test"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestRepository_GetGamesByIds(t *testing.T) {
	r := NewRepository(test.DB(), test.PubSub())

	res, err := r.GetGamesByIds([]string{"1", "2"})
	assert.NoError(t, err)

	assert.Len(t, res, 2)
	assert.Equal(t, models.GameStateNegotiation, res[0].State)
	assert.Equal(t, models.GameStateInProgress, res[1].State)
}

func TestRepository_GetGamesByState(t *testing.T) {
	r := NewRepository(test.DB(), test.PubSub())

	res, err := r.GetGamesByState([]models.GameState{models.GameStateNegotiation})
	assert.NoError(t, err)
	assert.True(t, len(res) > 0)
	assert.Len(t, res, 1)

	res, err = r.GetGamesByState([]models.GameState{models.GameStateNegotiation, models.GameStateInProgress})
	assert.NoError(t, err)
	assert.True(t, len(res) > 1)
}

func TestRepository_CreateGameUser(t *testing.T) {
	r := NewRepository(test.DB(), test.PubSub())

	res, err := r.CreateGameUser("1", "1", models.GameUserEdgeTypePlayer)
	assert.NoError(t, err)
	assert.Equal(t, "1", res.Id)
}

func TestRepository_CreateGameUserStateChange(t *testing.T) {
	r := NewRepository(test.DB(), test.PubSub())

	hash, err := bcrypt.GenerateFromPassword([]byte("hunter2"), 4)
	assert.NoError(t, err)

	id, err := r.CreateIdentity("testgameprovider_gameuserstatechange@tengen.io", hash, "test user")
	assert.NoError(t, err)

	user, err := r.GetUserById("1")
	assert.NoError(t, err)

	game, err := r.CreateGame(models.GameTypeStandard, 19, models.GameStateNegotiation, []models.User{id.User})
	assert.NoError(t, err)
	assert.Equal(t, game.State, models.GameStateNegotiation)

	game, err = r.CreateGameUser(game.Id, user.Id, models.GameUserEdgeTypePlayer)
	assert.NoError(t, err)
}

func TestRepository_CreateInvitation(t *testing.T) {
	r := NewRepository(test.DB(), test.PubSub())

	identity := models.Identity{
		User: models.User{
			NodeFields: models.NodeFields{
				Id: "1",
			},
		},
	}

	res, err := r.CreateGame(models.GameTypeStandard, 19, models.GameStateNegotiation, []models.User{identity.User})

	assert.NoError(t, err)
	assert.Equal(t, models.GameTypeStandard, res.Type)
	assert.Equal(t, models.GameStateNegotiation, res.State)
	assert.Equal(t, 19, res.BoardSize)
}

func TestRepository_GetIdentityById(t *testing.T) {
	r := NewRepository(test.DB(), test.PubSub())

	res, err := r.GetIdentityById(1)
	assert.NoError(t, err)

	assert.Equal(t, res.Email, "test1@tengen.io")
	assert.Equal(t, res.Name, "Test User 1")
}

func TestRepository_CreateIdentity(t *testing.T) {
	r := NewRepository(test.DB(), test.PubSub())

	hash, err := bcrypt.GenerateFromPassword([]byte("hunter2"), 4)
	assert.NoError(t, err)
	res, err := r.CreateIdentity("test-createidentity@tengen.io", hash, "Test User CreateIdentity")
	assert.NoError(t, err)

	assert.Equal(t, "test-createidentity@tengen.io", res.Email)
	assert.Equal(t, "Test User CreateIdentity", res.Name)
}

func TestRepository_GetUserById(t *testing.T) {
	r := NewRepository(test.DB(), test.PubSub())

	res, err := r.GetUserById("1")
	assert.NoError(t, err)

	assert.Equal(t, "1", res.Id)
	assert.Equal(t, "Test User 1", res.Name)
}

func TestMain(m *testing.M) {
	test.TestMain(m, "repository")
}
