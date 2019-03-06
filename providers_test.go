package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/tengen-io/server/models"
	"testing"
	"time"
)

func TestAuth_SignAndVerifyJWT(t *testing.T) {
	db := MakeTestDb()
	duration, _ := time.ParseDuration("1 week")
	auth := NewAuthProvider(db, []byte("supersecret"), duration)

	user := models.Identity{
		NodeFields: models.NodeFields{
			Id: "1",
		},
		Email: "test@test.com",
	}

	tokenStr, err := auth.SignJWT(user)
	assert.NoError(t, err)

	token, err := auth.ValidateJWT(tokenStr)
	assert.NoError(t, err)

	claims, ok := token.Claims.(*jwt.StandardClaims)
	assert.True(t, ok)
	assert.Equal(t, claims.Id, "1")
	assert.Equal(t, claims.Issuer, "tengen.io")
}

func TestAuth_ValidateInvalidJWT(t *testing.T) {
	db := MakeTestDb()
	duration, _ := time.ParseDuration("1 week")
	auth := NewAuthProvider(db, []byte("supersecret"), duration)

	_, err := auth.ValidateJWT("lol this wont work")
	assert.Error(t, err)
}

func TestGameProvider_GetGamesByIds(t *testing.T) {
	db := MakeTestDb()
	p := NewGameProvider(db)

	res, err := p.GetGamesByIds([]string{"1", "2"})
	assert.NoError(t, err)

	assert.Len(t, res, 2)
	assert.Equal(t, models.Invitation, res[0].State)
	assert.Equal(t, models.InProgress, res[1].State)
}

func TestGameProvider_GetGamesByState(t *testing.T) {
	db := MakeTestDb()
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
	db := MakeTestDb()
	p := NewGameProvider(db)

	res, err := p.CreateInvitation(models.CreateGameInvitationInput{
		BoardSize: 19,
		Type:      models.Standard,
	})

	assert.NoError(t, err)
	assert.Equal(t, models.Standard, res.Type)
	assert.Equal(t, models.Invitation, res.State)
	assert.Equal(t, 19, res.BoardSize)
}

func TestIdentityProvider_GetIdentityById(t *testing.T) {
	db := MakeTestDb()
	p := NewIdentityProvider(db, 4)
	res, err := p.GetIdentityById(1)
	assert.NoError(t, err)

	assert.Equal(t, res.Email, "test1@tengen.io")
	assert.Equal(t, res.Name, "Test User 1")
}

func TestIdentityProvider_CreateIdentity(t *testing.T) {
	db := MakeTestDb()
	p := NewIdentityProvider(db, 4)
	res, err := p.CreateIdentity("test-createidentity@tengen.io", "hunter2", "Test User CreateIdentity")
	assert.NoError(t, err)

	assert.Equal(t, "test-createidentity@tengen.io", res.Email)
	assert.Equal(t, "Test User CreateIdentity", res.Name)
}

func TestUserProvider_GetUserById(t *testing.T) {
	db := MakeTestDb()
	p := NewUserProvider(db)

	res, err := p.GetUserById("1")
	assert.NoError(t, err)

	assert.Equal(t, "1", res.Id)
	assert.Equal(t, "Test User 1", res.Name)
}
