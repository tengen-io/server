/*
Package resolvers is responsible for handling GraphQL queries and mutations.
*/
package resolvers

import (
	"github.com/camirmas/go_stop/models"
	"github.com/graphql-go/graphql"
)

// GetUser gets a User by username.
func GetUser(p graphql.ResolveParams) (interface{}, error) {
	db := p.Context.Value("db").(models.Database)
	return db.GetUser(p.Args["username"].(string))
}

// CreateUser creates a new User account.
func CreateUser(p graphql.ResolveParams) (interface{}, error) {
	db := p.Context.Value("db").(models.Database)
	username := p.Args["username"].(string)
	email := p.Args["email"].(string)
	password := p.Args["password"].(string)
	passwordConfirm := p.Args["passwordConfirmation"].(string)

	user, err := db.CreateUser(username, email, password, passwordConfirm)

	if err != nil {
		return nil, err
	}

	token, err := GenerateToken(user.Id)

	return &models.AuthUser{token, user}, nil
}

// LogIn generates a new JWT given a valid username/password.
func LogIn(p graphql.ResolveParams) (interface{}, error) {
	username := p.Args["username"].(string)
	password := p.Args["password"].(string)
	db := p.Context.Value("db").(models.Database)
	user, err := db.CheckPw(username, password)

	if err != nil {
		return nil, err
	}

	token, err := GenerateToken(user.Id)

	if err != nil {
		return user, err
	}

	authUser := &models.AuthUser{token, user}

	return authUser, nil
}

// CurrentUser gets the user corresponding to the provided JWT from request
// Content Headers.
func CurrentUser(p graphql.ResolveParams) (interface{}, error) {
	db := p.Context.Value("db").(models.Database)

	token, ok := p.Context.Value("token").(string)

	if !ok {
		return nil, missingTokenError{}
	}

	claims, err := ValidateToken(token)

	if err != nil {
		return nil, err
	}

	return db.GetUser(claims.UserId)
}
