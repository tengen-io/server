package server

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestHomepage(t *testing.T) {
	recorder := httptest.NewRecorder()
	config := ServerConfig{"", 0, false}
	server := Server{config: &config}
	handler := server.getHomepageHandler()

	request, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatalf("Unable to build a new request err: %s", err)
	}

	handler(recorder, request)
	response := recorder.Result()

	if response.StatusCode != 200 {
		t.Fatalf("Homepage did not respond successfully code: %d", response.StatusCode)
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Unable to read response body err: %s", err)
	}

	htmlTag, err := regexp.Match("<html>", responseBody)
	if err != nil {
		t.Fatalf("Unable to scan body for HTML tag err: %s", err)
	}

	if !htmlTag {
		t.Errorf("Homepage did not respond with valid HTML %s", responseBody)
	}
}
