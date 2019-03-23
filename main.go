package main

import (
	"github.com/joho/godotenv"
	"github.com/tengen-io/server/cmd"
	"os"
)

func main() {
	env := os.Getenv("TENGEN_ENV")
	if env == "" {
		env = "development"
	}

	godotenv.Load(".env." + env)
	godotenv.Load()

	cmd.Execute()
}
