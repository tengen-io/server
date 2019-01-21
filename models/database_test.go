package models

import (
	"log"
	"os"
	"testing"
)

var db *DB

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	teardown()

	os.Exit(code)
}

func TestCheckPw(t *testing.T) {
	db.CreateUser("testcheckpw", "testcheckpw@dude.dude", "dudedude", "dudedude")
	_, err := db.CheckPw("testcheckpw", "dudedude")

	if err != nil {
		t.Error(err)
	}
}

func expectErr(t *testing.T, expected, err error) {
	if err == nil {
		t.Errorf("Expected '%s'", expected.Error())
	}
	if err.Error() != expected.Error() {
		t.Errorf("Expected '%s', got '%s'", expected.Error(), err.Error())
	}
}

func setup() {
	newDb, err := ConnectDB()

	if err != nil {
		log.Fatal(err)
	}

	db = newDb
}

func teardown() {
	db.Query("DELETE FROM games")
	db.Query("DELETE FROM players")
	db.Query("DELETE FROM users")
}

func init() {
	testingMode = true
}
