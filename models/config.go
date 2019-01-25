package models

import (
	"fmt"
	"io/ioutil"
	"os"
)

const (
	testingDB string = "postgres://postgres:postgres@localhost:5432/go_stop_go_test?sslmode=disable"
	devDB     string = "postgres://postgres:postgres@localhost:5432/go_stop_go?sslmode=disable"

	testEnv string = "test"
	devEnv  string = "dev"
	prodEnv string = "prod"

	devBcryptRounds  int = 1
	prodBcryptRounds int = 12
)

type DbConfig struct {
	Env          string
	BcryptRounds int
	DbUrl        string
}

func (config *DbConfig) setEnv() {
	switch os.Getenv("ENV") {
	case testEnv:
		config.Env = testEnv
		config.BcryptRounds = devBcryptRounds
	case devEnv:
		config.Env = devEnv
		config.BcryptRounds = devBcryptRounds
	case prodEnv:
		config.Env = prodEnv
		config.BcryptRounds = prodBcryptRounds
	default:
		config.Env = devEnv
		config.BcryptRounds = devBcryptRounds
	}
}

func (config *DbConfig) setDbUrl() {
	var url string
	switch config.Env {
	case testEnv:
		url = os.Getenv("TEST_POSTGRES_URL")
		if url == "" {
			url = testingDB
		}
	case devEnv:
		url = os.Getenv("DEV_POSTGRES_URL")
		if url == "" {
			url = devDB
		}
	case prodEnv:
		dbName := os.Getenv("POSTGRES_DB")
		host := os.Getenv("POSTGRES_HOST")

		filename := os.Getenv("POSTGRES_USER")
		user := getSecret(filename)

		filename = os.Getenv("POSTGRES_PASSWORD")
		pw := getSecret(filename)

		url = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, pw, host, dbName)
	default:
		url = devDB
	}
	config.DbUrl = url
}

func getSecret(fileName string) string {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return ""
	}

	return string(file)
}

func setupConfig() *DbConfig {
	config := &DbConfig{}

	config.setEnv()
	config.setDbUrl()

	return config
}
