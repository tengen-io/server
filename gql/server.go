package gql

import (
	"encoding/base64"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/tengen-io/server/db"
	"github.com/tengen-io/server/pubsub"
	"github.com/tengen-io/server/repository"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type envConfig struct {
	Environment string
}

type serverConfig struct {
	Host            string
	Port            int
	GraphiQLEnabled bool
}

type server struct {
	config           *serverConfig
	executableSchema graphql.ExecutableSchema
	repo             *repository.Repository
	bcryptCost       int
	signingKey       []byte
	jwtLifetime      time.Duration
}

func (s *server) Start() {
	http.Handle("/graphql", enableCorsMiddleware(s.VerifyTokenMiddleware(handler.GraphQL(s.executableSchema))))
	http.Handle("/register", enableCorsMiddleware(s.RegistrationHandler()))
	http.Handle("/login", enableCorsMiddleware(s.LoginHandler()))
	http.HandleFunc("/", handler.Playground("tengen.io | GraphQL", "/graphql"))

	log.Printf("Listening on http://%s:%d", s.config.Host, s.config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", s.config.Host, s.config.Port), nil))
}

func enableCorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		next.ServeHTTP(w, r)
	})
}

func newServer(config *serverConfig, schema graphql.ExecutableSchema, repo *repository.Repository) *server {
	bcryptCost, err := strconv.Atoi(os.Getenv("TENGEN_BCRYPT_COST"))
	if err != nil {
		log.Fatal("Could not parse TENGEN_BCRYPT_COST")
	}

	day, err := time.ParseDuration("24h")
	if err != nil {
		log.Fatal("could not parse auth key duration", err)
	}

	keyDuration := day * 7

	return &server{
		config,
		schema,
		repo,
		bcryptCost,
		getSigningKey(),
		keyDuration,
	}
}

func makeRepo() *repository.Repository {
	port, err := strconv.Atoi(os.Getenv("TENGEN_DB_PORT"))
	if err != nil {
		log.Fatal("Cold not parse TENGEN_DB_PORT")
	}

	config := db.PostgresDBConfig{
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

	ps := pubsub.NewDbPubSub(config.Url())
	ps.Start()

	return repository.NewRepository(db, ps)
}

func makeSchema(repo *repository.Repository) graphql.ExecutableSchema {
	return NewExecutableSchema(Config{
		Resolvers: &Resolver{
			repo:   repo,
		},
		Directives: Directives(),
	})
}

func newServerConfig() envConfig {
	environment := os.Getenv("GO_ENV")

	if environment == "" {
		environment = "development"
	}

	return envConfig{
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

func makeServer(schema graphql.ExecutableSchema, repo *repository.Repository) *server {
	config := newServerConfig()
	tengenPort := os.Getenv("TENGEN_PORT")
	port, err := strconv.Atoi(tengenPort)
	if err != nil {
		log.Fatalf("Could not parse TENGEN_PORT err: %s", err)
	}

	serverConfig := &serverConfig{
		Host:            os.Getenv("TENGEN_HOST"),
		Port:            port,
		GraphiQLEnabled: config.Environment == "development",
	}

	return newServer(serverConfig, schema, repo)
}

func Serve() {
	repo := makeRepo()

	schema := makeSchema(repo)
	s := makeServer(schema, repo)
	s.Start()
}
