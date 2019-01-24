package resolvers

import (
	"context"
	"github.com/camirmas/go_stop/models"
	"github.com/graphql-go/graphql"
	"testing"
)

func TestCreateUser(t *testing.T) {
	db := setup()

	params := graphql.ResolveParams{}
	params.Args = map[string]interface{}{}
	params.Args["username"] = "dude"
	params.Args["email"] = "dude@dude.dude"
	params.Args["password"] = "dudedude"
	params.Args["passwordConfirmation"] = "dudedude"
	ctx := context.WithValue(params.Context, "db", db)
	params.Context = ctx

	_, err := CreateUser(params)

	if err != nil {
		t.Error(err)
	}
}

func setup() *models.TestDB {
	return &models.TestDB{}
}
