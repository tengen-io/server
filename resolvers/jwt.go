package resolvers

import (
	"github.com/dgrijalva/jwt-go"
)

type MyCustomClaims struct {
	UserId int `json:"userId"`
	jwt.StandardClaims
}

// GenerateToken generates a new JWT for the given User id.
func GenerateToken(userId int, signingKey []byte) (string, error) {
	// Create the Claims
	claims := MyCustomClaims{
		userId,
		jwt.StandardClaims{
			Issuer: "Tengen",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(signingKey)

	if err != nil {
		return "", err
	}
	return ss, nil
}

// ValidateToken parses and validates a provided token, retrieving the token's
// claims if successful.
func ValidateToken(tokenString string, signingKey []byte) (*MyCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	return token.Claims.(*MyCustomClaims), nil
}

type missingTokenError struct{}
type invalidTokenError struct{}

func (e missingTokenError) Error() string {
	return "Missing Authorization header"
}

func (e invalidTokenError) Error() string {
	return "Auth token invalid"
}
