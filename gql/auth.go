package gql

import (
	"context"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/handler"
	"github.com/dgrijalva/jwt-go"
	"github.com/tengen-io/server/models"
	"github.com/tengen-io/server/repository"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

type auth struct {
	signingKey []byte
	jwtLifetime time.Duration
	repo repository.Repository
	bcryptCost int
}

func (a *auth) validateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}

		return a.signingKey, nil
	})

	return token, err
}

func (a *auth) signJWT(identity models.Identity) (string, error) {
	// TODO(eac): reintroduce custom claims for ID
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Id:        identity.Id,
		NotBefore: time.Now().Unix(),
		ExpiresAt: time.Now().Add(a.jwtLifetime * time.Second).Unix(),
		Issuer:    "tengen.io",
	})

	ss, err := token.SignedString(a.signingKey)
	if err != nil {
		return "", err
	}

	return ss, nil
}

func (a *auth) checkPasswordByEmail(email, password string) (*models.Identity, error) {
	var rv *models.Identity

	err := a.repo.WithTx(func(r *repository.Repository) error {
		passwordHash, err := a.repo.GetPwHashForEmail(email)
		if err != nil {
			return err
		}

		passwordBytes := []byte(password)
		err = bcrypt.CompareHashAndPassword(passwordHash, passwordBytes)
		if err != nil {
			return err
		}

		rv, err = a.repo.GetIdentityByEmail(email)
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

func (a *auth) authForContext(ctx context.Context) (*models.Identity, error) {
	id, ok := ctx.Value(IdentityContextKey).(models.Identity)
	if !ok {
		// Could be websocket
		initPayload := handler.GetInitPayload(ctx)
		tokenStr := initPayload.Authorization()
		if tokenStr == "" {
			return nil, errors.New("unauthorized")
		}

		token, err := a.validateJWT(tokenStr)
		if err != nil {
			return nil, err
		}

		claims, ok := token.Claims.(*jwt.StandardClaims)
		if !ok {
			return nil, errors.New("unable to cast claims")
		}

		idInt, err := strconv.Atoi(claims.Id)
		if err != nil {
			return nil, err
		}

		identity, err := a.repo.GetIdentityById(int32(idInt))
		if err != nil {
			return nil, err
		}

		return identity, nil
	}

	return &id, nil
}
