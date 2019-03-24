package gql

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"github.com/tengen-io/server/models"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthRepository struct {
	db         *sqlx.DB
	signingKey []byte
	lifetime   time.Duration
}

func NewAuthRepository(db *sqlx.DB, signingKey []byte, lifetime time.Duration) *AuthRepository {
	return &AuthRepository{
		db:         db,
		signingKey: signingKey,
		lifetime:   lifetime,
	}
}

func (p *AuthRepository) ValidateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}

		return p.signingKey, nil
	})

	return token, err
}

func (p *AuthRepository) SignJWT(identity models.Identity) (string, error) {
	// TODO(eac): reintroduce custom claims for ID
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Id:        identity.Id,
		NotBefore: time.Now().Unix(),
		ExpiresAt: time.Now().Add(p.lifetime * time.Second).Unix(),
		Issuer:    "tengen.io",
	})

	ss, err := token.SignedString(p.signingKey)
	if err != nil {
		return "", err
	}

	return ss, nil
}

// TODO(eac): Figure out how to use dbx structs for nested structures
func (p *AuthRepository) CheckPasswordByEmail(email, password string) (*models.Identity, error) {
	var passwordHash string
	err := p.db.Get(&passwordHash, "SELECT password_hash FROM identities WHERE email = $1", email)
	if err != nil {
		return nil, err
	}

	passwordBytes := []byte(password)
	hashBytes := []byte(passwordHash)

	err = bcrypt.CompareHashAndPassword(hashBytes, passwordBytes)
	if err != nil {
		return nil, err
	}

	var rv models.Identity
	row := p.db.QueryRowx("SELECT i.id, i.email, u.id, u.name FROM identities i, users u WHERE i.id = u.identity_id AND email = $1", email)
	err = row.Scan(&rv.Id, &rv.Email, &rv.User.Id, &rv.User.Name)

	if err != nil {
		return nil, err
	}

	return &rv, nil
}
