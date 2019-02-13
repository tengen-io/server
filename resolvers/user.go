/*
Package resolvers is responsible for handling GraphQL queries and mutations.
*/
package resolvers

import (
	"github.com/camirmas/go_stop/models"
	"github.com/graphql-go/graphql"
)

// GetUser gets a User by username.
func (r *Resolvers) GetUser(p graphql.ResolveParams) (interface{}, error) {
	return r.db.GetUser(p.Args["username"].(string))
}

// CreateUser creates a new User account.
func (r *Resolvers) CreateUser(p graphql.ResolveParams) (interface{}, error) {
	username := p.Args["username"].(string)
	email := p.Args["email"].(string)
	password := p.Args["password"].(string)
	passwordConfirm := p.Args["passwordConfirmation"].(string)

	user, err := r.db.CreateUser(username, email, password, passwordConfirm)

	if err != nil {
		return nil, err
	}

	token, err := GenerateToken(user.Id, r.signingKey)

	return &models.AuthUser{token, user}, nil
}

// LogIn generates a new JWT given a valid username/password.
func (r *Resolvers) LogIn(p graphql.ResolveParams) (interface{}, error) {
	username := p.Args["username"].(string)
	password := p.Args["password"].(string)
	user, err := r.db.CheckPw(username, password)

	if err != nil {
		return nil, err
	}

	token, err := GenerateToken(user.Id, r.signingKey)

	if err != nil {
		return user, err
	}

	authUser := &models.AuthUser{token, user}

	return authUser, nil
}

// CurrentUser gets the user corresponding to the provided JWT from request
// Content Headers.
func (r *Resolvers) CurrentUser(p graphql.ResolveParams) (interface{}, error) {
	token, ok := p.Context.Value("token").(string)

	if !ok {
		return nil, missingTokenError{}
	}

	claims, err := ValidateToken(token, r.signingKey)

	if err != nil {
		return nil, err
	}

	return r.db.GetUser(claims.UserId)
}
