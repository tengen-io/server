package server

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/tengen-io/server/providers"
	"github.com/tengen-io/server/test"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestServer_LoginHandler(t *testing.T) {
	server := makeServer()
	reqBody := "{\"email\": \"test1@tengen.io\", \"password\": \"hunter2\"}"

	req, err := http.NewRequest("POST", "/login", strings.NewReader(reqBody))
	assert.NoError(t, err)

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

func makeServer() *Server {
	db := test.MakeDb()
	config := ServerConfig{
		"", 0, false,
	}
	duration, _ := time.ParseDuration("1 week")
	identityProvider := providers.NewIdentityProvider(db, 1)
	return NewServer(&config, nil, providers.NewAuthProvider(db, []byte("supersecret"), duration), identityProvider)
}


func TestMain(m *testing.M) {
	test.Main(m)
}
