package main

import (
	"crypto/sha256"
	"github.com/dgrijalva/jwt-go"
)

func createToken(username string) (string, error) {
	type MyCustomClaims struct {
		Username string `json:"username"`
		jwt.StandardClaims
	}

	h := sha256.New()
	signingKey := h.Sum([]byte("TODO: pull secret key from safe place"))

	// Create the Claims
	claims := MyCustomClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    "GoStop",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(signingKey)

	if err != nil {
		return "", err
	}
	return ss, nil
}
