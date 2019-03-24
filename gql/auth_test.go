package gql

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/tengen-io/server/models"
	"testing"
)

func TestServer_SignAndVerifyJWT(t *testing.T) {
	server := makeTestServer()

	user := models.Identity{
		NodeFields: models.NodeFields{
			Id: "1",
		},
		Email: "test@test.com",
	}

	tokenStr, err := server.signJWT(user)
	assert.NoError(t, err)

	token, err := server.validateJWT(tokenStr)
	assert.NoError(t, err)

	claims, ok := token.Claims.(*jwt.StandardClaims)
	assert.True(t, ok)
	assert.Equal(t, claims.Id, "1")
	assert.Equal(t, claims.Issuer, "tengen.io")
}

func TestServer_ValidateInvalidJWT(t *testing.T) {
	server := makeTestServer()
	_, err := server.validateJWT("lol this wont work")
	assert.Error(t, err)
}
