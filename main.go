package main

import (
	"encoding/base64"
	"github.com/camirmas/go_stop/models"
	"github.com/camirmas/go_stop/providers"
	"github.com/camirmas/go_stop/resolvers"
	"github.com/camirmas/go_stop/server"
	"github.com/graphql-go/graphql"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

func getSigningKey() []byte {
	encoded := os.Getenv("GOSTOP_JWT_SECRET_KEY")
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		log.Fatal("Could not load JWT signing key.", err)
	}

	return decoded
}

func makeServer(db models.DB, auth *providers.Auth, schema *graphql.Schema) *server.Server {
	port, err := strconv.Atoi(os.Getenv("GOSTOP_PORT"))
	if err != nil {
		log.Fatal("Could not parse GOSTOP_PORT")
	}

	config := &server.ServerConfig{
		Host:       os.Getenv("GOSTOP_HOST"),
		Port:       port,
	}

	return server.NewServer(config, db, auth, schema)
}

func makeDb() *models.PostgresDB {
	port, err := strconv.Atoi(os.Getenv("GOSTOP_DB_PORT"))
	if err != nil {
		log.Fatal("Cold not parse GOSTOP_DB_PORT")
	}

	bcryptRounds, err := strconv.Atoi(os.Getenv("GOSTOP_BCRYPT_ROUNDS"))
	if err != nil {
		log.Fatal("Could not parse GOSTOP_BCRYPT_ROUNDS")
	}

	config := &models.PostgresDBConfig{
		Host:         os.Getenv("GOSTOP_DB_HOST"),
		Port:         port,
		User:         os.Getenv("GOSTOP_DB_USER"),
		Database:     os.Getenv("GOSTOP_DB_DATABASE"),
		Password:     os.Getenv("GOSTOP_DB_PASSWORD"),
		BcryptRounds: bcryptRounds,
	}

	db, err := models.NewPostgresDB(config)
	if err != nil {
		log.Fatal("Unable to connect to DB.", err)
	}

	return db
}

func makeAuth() *providers.Auth {
	keyDuration, err := time.ParseDuration("1 week")
	if err != nil {
		log.Fatal("could not parse auth key duration", err)
	}

	return providers.NewAuth(getSigningKey(), keyDuration)
}

func makeResolvers(db models.DB, auth *providers.Auth) *resolvers.Resolvers {
	return resolvers.NewResolvers(db, auth)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Could not load .env", err)
	}

	db := makeDb()
	auth := makeAuth()
	resolvers := makeResolvers(db, auth)

	schema, err := NewSchema(resolvers)
	if err != nil {
		log.Fatal("Could not crete GraphQL schema", err)
	}

	s := makeServer(db, auth, &schema)
	s.Start()
}
