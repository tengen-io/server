package main

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestServer_LoginHandler(t *testing.T) {
	server := makeTestServer()
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

func makeTestServer() *Server {
	db := MakeTestDb()
	config := ServerConfig{
		"", 0, false,
	}
	duration, _ := time.ParseDuration("1 week")
	identityProvider := NewIdentityProvider(db, 1)
	return NewServer(&config, nil, NewAuthProvider(db, []byte("supersecret"), duration), identityProvider)
}
