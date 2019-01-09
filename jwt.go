package main

import (
	"crypto/sha256"
	"github.com/dgrijalva/jwt-go"
)

type MyCustomClaims struct {
	UserId int `json:"userId"`
	jwt.StandardClaims
}

func getKey() []byte {
	h := sha256.New()

	return h.Sum([]byte("TODO: pull secret key from safe place"))
}

func GenerateToken(userId int) (string, error) {
	// Create the Claims
	claims := MyCustomClaims{
		userId,
		jwt.StandardClaims{
			Issuer: "GoStop",
		},
	}
	signingKey := getKey()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(signingKey)

	if err != nil {
		return "", err
	}
	return ss, nil
}

func ValidateToken(tokenString string) (*MyCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return getKey(), nil
	})

	if err != nil {
		return nil, err
	}

	return token.Claims.(*MyCustomClaims), nil
}
