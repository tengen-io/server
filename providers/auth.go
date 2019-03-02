package providers

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"github.com/tengen-io/server/models"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

type AuthProvider struct {
	db         *sqlx.DB
	signingKey []byte
	lifetime   time.Duration
}

func NewAuthProvider(db *sqlx.DB, signingKey []byte, lifetime time.Duration) *AuthProvider {
	return &AuthProvider{
		signingKey: signingKey,
		lifetime:   lifetime,
	}
}

func (p *AuthProvider) ValidateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}

		return p.signingKey, nil
	})

	return token, err
}

func (p *AuthProvider) SignJWT(identity models.Identity) (string, error) {
	// TODO(eac): reintroduce custom claims for ID
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Id:        strconv.Itoa(int(identity.Id)),
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

func (p *AuthProvider) CheckPasswordByEmail(email, password string) (*models.Identity, error) {
	var passwordHash string
	err := p.db.Select(&passwordHash, "SELECT passwordHash FROM identities WHERE email = ?", email)
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
	err = p.db.Select(&rv, "SELECT i.*, u.* FROM identities i JOIN users u USING (i.id, u.identityId) WHERE email = ?", email)
	if err != nil {
		return nil, err
	}

	return &rv, nil
}
