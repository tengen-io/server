/*
Server implementation of the board game Go.
*/
package server

import (
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/handler"
	"github.com/tengen-io/server/providers"
)

type ServerConfig struct {
	Host            string
	Port            int
	GraphiQLEnabled bool
}

type Server struct {
	config           *ServerConfig
	executableSchema graphql.ExecutableSchema
	auth             *providers.AuthProvider
	identity         *providers.IdentityProvider
}

func NewServer(config *ServerConfig, schema graphql.ExecutableSchema, auth *providers.AuthProvider, identity *providers.IdentityProvider) *Server {
	return &Server{
		config,
		schema,
		auth,
		identity,
	}
}

func (s *Server) Start() {
	http.Handle("/graphql", enableCorsMiddleware(s.VerifyTokenMiddleware(handler.GraphQL(s.executableSchema))))
	http.Handle("/login", s.LoginHandler())
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
