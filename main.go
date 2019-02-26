package main

import (
	"encoding/base64"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/camirmas/go_stop/models"
	"github.com/camirmas/go_stop/providers"
	"github.com/camirmas/go_stop/resolvers"
	"github.com/camirmas/go_stop/server"
	"github.com/graphql-go/graphql"
	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	Port        int
}

func NewConfig() Config {
	environment := os.Getenv("GO_ENV")
	goStopPort := os.Getenv("GOSTOP_PORT")

	if environment == "" {
		environment = "development"
	}

	port, err := strconv.Atoi(goStopPort)
	if err != nil {
		log.Fatalf("Could not parse GOSTOP_PORT err: %s", err)
	}

	return Config{
		Environment: environment,
		Port:        port,
	}
}

func getSigningKey() []byte {
	encoded := os.Getenv("GOSTOP_JWT_SECRET_KEY")
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		log.Fatal("Could not load JWT signing key.", err)
	}

	return decoded
}

func makeServer(db models.DB, auth *providers.Auth, schema *graphql.Schema) *server.Server {
	config := NewConfig()

	serverConfig := &server.ServerConfig{
		Host:            os.Getenv("GOSTOP_HOST"),
		Port:            config.Port,
		GraphiQLEnabled: config.Environment == "development",
	}

	return server.NewServer(serverConfig, db, auth, schema)
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
	day, err := time.ParseDuration("24h")
	if err != nil {
		log.Fatal("could not parse auth key duration", err)
	}

	keyDuration := day * 7

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
