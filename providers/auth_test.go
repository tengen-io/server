package providers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/tengen-io/server/models"
	"testing"
	"time"
)

func TestAuth_SignAndVerifyJWT(t *testing.T) {
	duration, _ := time.ParseDuration("1 week")
	auth := NewAuthProvider([]byte("supersecret"), duration)

	user := models.User{
		Id:       1,
		Username: "test",
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
	duration, _ := time.ParseDuration("1 week")
	auth := NewAuthProvider([]byte("supersecret"), duration)

	_, err := auth.ValidateJWT("lol this wont work")
	assert.Error(t, err)
}
