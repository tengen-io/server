package providers

import (
	"github.com/camirmas/go_stop/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAuth_SignAndVerifyJWT(t *testing.T) {
	duration, _ := time.ParseDuration("1 week")
	auth := NewAuth([]byte("supersecret"), duration)

	user := models.User{
		Id: 1,
		Username: "test",
		Email: "test@test.com",
	}

	tokenStr, err := auth.SignJWT(user)
	assert.NoError(t, err)

	token, err := auth.ValidateJWT(tokenStr)
	assert.NoError(t, err)

	claims, ok := token.Claims.(*jwt.StandardClaims)
	assert.True(t, ok)
	assert.Equal(t, claims.Id, "1")
	assert.Equal(t, claims.Issuer, "gostop")
}

func TestAuth_ValidateInvalidJWT(t *testing.T) {
	duration, _ := time.ParseDuration("1 week")
	auth := NewAuth([]byte("supersecret"), duration)

	_, err := auth.ValidateJWT("lol this wont work")
	assert.Error(t, err)
}