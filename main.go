package main

import (
	"encoding/base64"
	"github.com/99designs/gqlgen/graphql"
	"github.com/jmoiron/sqlx"
	"github.com/tengen-io/server/db"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type envConfig struct {
	Environment string
}

func newServerConfig() envConfig {
	environment := os.Getenv("GO_ENV")

	if environment == "" {
		environment = "development"
	}

	return envConfig {
		Environment: environment,
	}
}

func getSigningKey() []byte {
	encoded := os.Getenv("TENGEN_JWT_SECRET_KEY")
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		log.Fatal("Could not load JWT signing key.", err)
	}

	return decoded
}

func makeServer(schema graphql.ExecutableSchema, auth *AuthProvider, identity *IdentityProvider) *Server {
	config := newServerConfig()
	tengenPort := os.Getenv("TENGEN_PORT")
	port, err := strconv.Atoi(tengenPort)
	if err != nil {
		log.Fatalf("Could not parse TENGEN_PORT err: %s", err)
	}

	serverConfig := &ServerConfig{
		Host:            os.Getenv("TENGEN_HOST"),
		Port:            port,
		GraphiQLEnabled: config.Environment == "development",
	}

	return NewServer(serverConfig, schema, auth, identity)
}

func makeDb() *sqlx.DB {
	port, err := strconv.Atoi(os.Getenv("TENGEN_DB_PORT"))
	if err != nil {
		log.Fatal("Cold not parse TENGEN_DB_PORT")
	}

	config := &db.PostgresDBConfig{
		Host:     os.Getenv("TENGEN_DB_HOST"),
		Port:     port,
		User:     os.Getenv("TENGEN_DB_USER"),
		Database: os.Getenv("TENGEN_DB_DATABASE"),
		Password: os.Getenv("TENGEN_DB_PASSWORD"),
	}

	db, err := db.NewPostgresDb(config)
	if err != nil {
		log.Fatal("Unable to connect to DB.", err)
	}

	return db
}

func makeAuth(db *sqlx.DB) *AuthProvider {
	day, err := time.ParseDuration("24h")
	if err != nil {
		log.Fatal("could not parse auth key duration", err)
	}

	keyDuration := day * 7

	return NewAuthProvider(db, getSigningKey(), keyDuration)
}

func makeIdentity(db *sqlx.DB) *IdentityProvider {
	bcryptCost, err := strconv.Atoi(os.Getenv("TENGEN_BCRYPT_COST"))
	if err != nil {
		log.Fatal("Could not parse TENGEN_BCRYPT_COST")
	}

	return NewIdentityProvider(db, bcryptCost)
}

func makeSchema(identity *IdentityProvider, user *UserProvider, game *GameProvider) graphql.ExecutableSchema {
	return NewExecutableSchema(Config{
		Resolvers: &Resolver{
			identity: identity,
			user:     user,
			game:     game,
		},
		Directives: Directives(),
	})
}

func main() {
	env := os.Getenv("TENGEN_ENV")
	if env == "" {
		env = "development"
	}

	godotenv.Load(".env." + env)
	godotenv.Load()

	db := makeDb()
	auth := makeAuth(db)
	identity := makeIdentity(db)
	user := NewUserProvider(db)
	game := NewGameProvider(db)
	schema := makeSchema(identity, user, game)

	s := makeServer(schema, auth, identity)
	s.Start()
}
