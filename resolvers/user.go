/*
Package resolvers is responsible for handling GraphQL queries and mutations.
*/
package resolvers

import (
	"github.com/graphql-go/graphql"
	"github.com/tengen-io/server/models"
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

	token, err := r.auth.SignJWT(*user)
	if err != nil {
		return nil, err
	}

	return &models.AuthUser{token, user}, nil
}

// CurrentUser gets the user corresponding to the provided JWT from request
// Content Headers.
func (r *Resolvers) CurrentUser(p graphql.ResolveParams) (interface{}, error) {
	currentUser, ok := p.Context.Value("currentUser").(*models.User)

	if !ok {
		return nil, currentUserError{}
	}

	return currentUser, nil
}
