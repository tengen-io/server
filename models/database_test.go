package models

import (
	"fmt"
	"log"
	"os"
	"testing"
)

var db *PostgresDB

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	teardown()

	os.Exit(code)
}

func TestCheckPw(t *testing.T) {
	_, err := db.CreateUser("testcheckpw", "testcheckpw@dude.dude", "dudedude", "dudedude")
	if err != nil {
		t.Error(err)
	}

	_, err = db.CheckPw("testcheckpw", "dudedude")

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
	config := &PostgresDBConfig{
		"localhost",
		5432,
		"postgres",
		"postgres",
		"tengen_test",
		1,
	}
	newDb, err := NewPostgresDB(config)

	if err != nil {
		log.Fatalf("Could not connect to test Postgres database: %s", err)
	}

	db = newDb
}

func teardown() {
	_, err := db.Exec("DELETE FROM players")
	if err != nil {
		fmt.Println(err)
	}
	_, err = db.Exec("DELETE FROM stones")
	if err != nil {
		fmt.Println(err)
	}
	_, err = db.Exec("DELETE FROM games")
	if err != nil {
		fmt.Println(err)
	}
	_, err = db.Exec("DELETE FROM users")
	if err != nil {
		fmt.Println(err)
	}
}
