package providers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/tengen-io/server/models"
	"github.com/tengen-io/server/test"
	"testing"
	"time"
)

func TestAuth_SignAndVerifyJWT(t *testing.T) {
	db := test.MakeDb()
	duration, _ := time.ParseDuration("1 week")
	auth := NewAuthProvider(db, []byte("supersecret"), duration)

	user := models.Identity{
		NodeFields: models.NodeFields {
			Id: "1",
		},
		Email:    "test@test.com",
	}

	tokenStr, err := auth.SignJWT(user)
	assert.NoError(t, err)

	token, err := auth.ValidateJWT(tokenStr)
	assert.NoError(t, err)

	claims, ok := token.Claims.(*jwt.StandardClaims)
	assert.True(t, ok)
	assert.Equal(t, claims.Id, "1")
	assert.Equal(t, claims.Issuer, "tengen")
}

func TestAuth_ValidateInvalidJWT(t *testing.T) {
	db := test.MakeDb()
	duration, _ := time.ParseDuration("1 week")
	auth := NewAuthProvider(db, []byte("supersecret"), duration)

	_, err := auth.ValidateJWT("lol this wont work")
	assert.Error(t, err)
}
