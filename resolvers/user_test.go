package resolvers

import (
	"context"
	"github.com/camirmas/go_stop/models"
	"github.com/graphql-go/graphql"
	"testing"
)

func TestCreateUser(t *testing.T) {
	r, params := setup()

	params.Args["username"] = "dude"
	params.Args["email"] = "dude@dude.dude"
	params.Args["password"] = "dudedude"
	params.Args["passwordConfirmation"] = "dudedude"

	u, err := r.CreateUser(params)

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
	r, params := setup()
	params.Args["username"] = "dude"

	_, err := r.GetUser(params)

	if err != nil {
		t.Error("Expected User, got error")
	}
}

func TestLogIn(t *testing.T) {
	r, params := setup()
	params.Args["username"] = "dude"
	params.Args["password"] = "dudedude"

	user, err := r.LogIn(params)

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
	r, params := setupAuth()
	_, err := r.CurrentUser(params)

	if err != nil {
		t.Error("Expected User, got error")
	}
}

func setup() (*Resolvers, graphql.ResolveParams) {
	db := &models.TestDB{}
	signingKey := []byte("secret")
	params := graphql.ResolveParams{}
	params.Args = map[string]interface{}{}
	params.Context = context.Background()

	r := &Resolvers{
		db:         db,
		signingKey: signingKey,
	}

	return r, params
}

func setupAuth() (*Resolvers, graphql.ResolveParams) {
	r, params := setup()
	token, _ := GenerateToken(1, r.signingKey)
	params.Context = context.WithValue(params.Context, "token", token)

	return r, params
}

func expectErr(t *testing.T, expected, err error) {
	if err == nil {
		t.Errorf("Expected '%s'", expected.Error())
	}
	if err.Error() != expected.Error() {
		t.Errorf("Expected '%s', got '%s'", expected.Error(), err.Error())
	}
}
