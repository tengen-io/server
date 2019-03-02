package server

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tengen-io/server/models"
	"github.com/tengen-io/server/providers"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type mockDb struct {
	mock.Mock
}

func (m mockDb) CheckPw(username, password string) (*models.User, error) {
	args := m.Called(username, password)
	return args.Get(0).(*models.User), args.Error(1)
}

func (mockDb) GetUser(id interface{}) (*models.User, error) {
	panic("implement me")
}

func (mockDb) CreateUser(username, email, password, passwordConfirm string) (*models.User, error) {
	panic("implement me")
}

func (mockDb) GetGame(gameId interface{}) (*models.Game, error) {
	panic("implement me")
}

func (mockDb) GetGames(userId interface{}) ([]*models.Game, error) {
	panic("implement me")
}

func (mockDb) CreateGame(userId int, opponent *models.User) (*models.Game, error) {
	panic("implement me")
}

func (mockDb) UpdateGame(userId int, game *models.Game, toAdd models.Stone, toRemove []models.Stone) error {
	panic("implement me")
}

func (mockDb) Pass(userId int, game *models.Game) (*models.Game, error) {
	panic("implement me")
}

func TestServer_LoginHandler(t *testing.T) {
	db := &mockDb{}
	server := makeServer(db)
	user := models.User{Id: 1, Username: "asdf", Email: "asdf@derp.com"}
	reqBody := "{\"username\": \"asdf@derp.com\", \"password\": \"secretpass\"}"

	req, err := http.NewRequest("POST", "/login", strings.NewReader(reqBody))
	assert.NoError(t, err)

	db.On("CheckPw", "asdf@derp.com", "secretpass").Return(&user, nil)

	rr := httptest.NewRecorder()
	handler := server.LoginHandler()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, 200)

	var resp struct{ Token string }

	bodyBuf, _ := ioutil.ReadAll(rr.Body)
	json.Unmarshal(bodyBuf, &resp)

	_, err = jwt.Parse(resp.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte("supersecret"), nil
	})

	assert.NoError(t, err)
}

func makeServer(db models.DB) *Server {
	config := ServerConfig{
		"", 0, false,
	}
	duration, _ := time.ParseDuration("1 week")
	return NewServer(&config, db, providers.NewAuthProvider([]byte("supersecret"), duration), &graphql.Schema{})
}
