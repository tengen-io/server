package gql

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/tengen-io/server/models"
	"github.com/tengen-io/server/repository"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (s *server) validateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}

		return s.signingKey, nil
	})

	return token, err
}

func (s *server) signJWT(identity models.Identity) (string, error) {
	// TODO(eac): reintroduce custom claims for ID
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Id:        identity.Id,
		NotBefore: time.Now().Unix(),
		ExpiresAt: time.Now().Add(s.jwtLifetime * time.Second).Unix(),
		Issuer:    "tengen.io",
	})

	ss, err := token.SignedString(s.signingKey)
	if err != nil {
		return "", err
	}

	return ss, nil
}

func (s *server) checkPasswordByEmail(email, password string) (*models.Identity, error) {
	var rv *models.Identity

	err := s.repo.WithTx(func(r *repository.Repository) error {
		passwordHash, err := s.repo.GetPwHashForEmail(email)
		if err != nil {
			return err
		}

		passwordBytes := []byte(password)
		err = bcrypt.CompareHashAndPassword(passwordHash, passwordBytes)
		if err != nil {
			return err
		}

		rv, err = s.repo.GetIdentityByEmail(email)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return rv, nil
}
