package resolvers

import (
	"context"
	"github.com/camirmas/go_stop/models"
	"github.com/graphql-go/graphql"
	"testing"
)

func TestCreateUser(t *testing.T) {
	params := setup()

	params.Args["username"] = "dude"
	params.Args["email"] = "dude@dude.dude"
	params.Args["password"] = "dudedude"
	params.Args["passwordConfirmation"] = "dudedude"

	u, err := CreateUser(params)

	if err != nil {
		t.Error(err)
	}

	user, ok := u.(*models.AuthUser)

	if !ok {
		t.Error("Expected AuthUser result")
	}

	if user.Jwt == "" {
		t.Error("Expected JWT to exist for User")
	}
}

func TestGetUser(t *testing.T) {
	params := setup()
	params.Args["username"] = "dude"

	_, err := GetUser(params)

	if err != nil {
		t.Error("Expected User, got error")
	}
}

func TestLogIn(t *testing.T) {
	params := setup()
	params.Args["username"] = "dude"
	params.Args["password"] = "dudedude"

	user, err := LogIn(params)

	if err != nil {
		t.Error("Expected logged in User")
	}

	authUser, ok := user.(*models.AuthUser)

	if !ok {
		t.Error("Expected AuthUser")
	}

	if authUser.Jwt == "" {
		t.Error("Expected JWT to exist for User")
	}
}

func TestCurrentUser(t *testing.T) {
	db := &models.TestDB{}
	params := graphql.ResolveParams{}
	params.Args = map[string]interface{}{}
	ctx := context.WithValue(params.Context, "db", db)
	token, _ := GenerateToken(1)
	ctx = context.WithValue(ctx, "token", token)
	params.Context = ctx

	_, err := CurrentUser(params)

	if err != nil {
		t.Error("Expected User, got error")
	}
}

func setup() graphql.ResolveParams {
	db := &models.TestDB{}
	params := graphql.ResolveParams{}
	params.Args = map[string]interface{}{}
	ctx := context.WithValue(params.Context, "db", db)
	params.Context = ctx

	return params
}

func setupAuth() graphql.ResolveParams {
	db := &models.TestDB{}
	params := graphql.ResolveParams{}
	params.Args = map[string]interface{}{}
	ctx := context.WithValue(params.Context, "db", db)
	token, _ := GenerateToken(1)
	ctx = context.WithValue(ctx, "token", token)
	params.Context = ctx

	return params
}

func expectErr(t *testing.T, expected, err error) {
	if err == nil {
		t.Errorf("Expected '%s'", expected.Error())
	}
	if err.Error() != expected.Error() {
		t.Errorf("Expected '%s', got '%s'", expected.Error(), err.Error())
	}
}
