package providers

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/tengen-io/server/models"
	"strconv"
	"time"
)

type Auth struct {
	signingKey []byte
	lifetime   time.Duration
}

func NewAuth(signingKey []byte, lifetime time.Duration) *Auth {
	return &Auth{
		signingKey: signingKey,
		lifetime:   lifetime,
	}
}

func (a *Auth) ValidateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}

		return a.signingKey, nil
	})

	return token, err
}

func (a *Auth) SignJWT(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Id:        strconv.Itoa(user.Id),
		NotBefore: time.Now().Unix(),
		ExpiresAt: time.Now().Add(a.lifetime * time.Second).Unix(),
		Issuer:    "tengen",
	})

	ss, err := token.SignedString(a.signingKey)
	if err != nil {
		return "", err
	}

	return ss, nil
}
