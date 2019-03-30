package gql

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/tengen-io/server/repository"
	"github.com/tengen-io/server/test"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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

func TestServer_RegistrationHandler(t *testing.T) {
	server := makeTestServer()

	testCases := []struct {
		testname         string
		email            string
		password         string
		username         string
		expectedResponse int
	}{
		{"add user", "test2@tengen.io", "test2", "test 2 user", 200},
		{"add empty user", "", "", "", 400},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testname, func(t *testing.T) {
			reqBody := fmt.Sprintf("{\"email\":\"%s\",\"password\":\"%s\",\"name\":\"%s\"}", testCase.email, testCase.password, testCase.username)
			req, err := http.NewRequest("POST", "/register", strings.NewReader(reqBody))
			assert.NoError(t, err)
			rr := httptest.NewRecorder()
			handler := server.RegistrationHandler()
			handler.ServeHTTP(rr, req)
			assert.Equal(t, testCase.expectedResponse, rr.Code)
		})
	}
}

func makeTestServer() *server {
	config := serverConfig{
		"", 0, false,
	}
	repo := repository.NewRepository(test.DB(), test.PubSub())
	return newServer(&config, nil, repo)
}

func TestMain(m *testing.M) {
	test.TestMain(m, "gql")
}
