package resolvers

import (
	"context"
	"github.com/graphql-go/graphql"
	"github.com/tengen-io/server/models"
	"github.com/tengen-io/server/providers"
	"testing"
	"time"
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

func TestCurrentUser(t *testing.T) {
	r, params := setupAuth()
	_, err := r.CurrentUser(params)

	if err != nil {
		t.Error("Expected User, got error")
	}
}

func setup() (*Resolvers, graphql.ResolveParams) {
	db := &models.TestDB{}
	params := graphql.ResolveParams{}
	params.Args = map[string]interface{}{}
	params.Context = context.Background()

	authDuration, _ := time.ParseDuration("1 week")
	auth := providers.NewAuthProvider([]byte("supersecret"), authDuration)

	r := &Resolvers{
		db:   db,
		auth: auth,
	}

	return r, params
}

func setupAuth() (*Resolvers, graphql.ResolveParams) {
	r, params := setup()
	user := models.User{Id: 1, Email: "test@example.org", Username: "test"}
	params.Context = context.WithValue(params.Context, "currentUser", &user)
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
