/*
Server implementation of the board game Go.
*/
package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/tengen-io/server/models"
	"github.com/tengen-io/server/providers"
)

type ServerConfig struct {
	Host            string
	Port            int
	GraphiQLEnabled bool
}

type Server struct {
	config *ServerConfig
	db     models.DB
	schema *graphql.Schema
	auth   *providers.Auth
}

func NewServer(config *ServerConfig, db models.DB, auth *providers.Auth, schema *graphql.Schema) *Server {
	return &Server{
		config,
		db,
		schema,
		auth,
	}
}

func (s *Server) Start() {
	config := s.config
	h := handler.New(&handler.Config{
		Schema:   s.schema,
		Pretty:   true,
		GraphiQL: config.GraphiQLEnabled,
	})

	http.Handle("/graphql", enableCorsMiddleware(s.VerifyTokenMiddleware(gqlMiddleware(h))))
	http.Handle("/login", s.LoginHandler())
	http.HandleFunc("/", s.getHomepageHandler())

	log.Printf("Listening on http://%s:%d", s.config.Host, s.config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", s.config.Host, s.config.Port), nil))
}

func gqlMiddleware(next *handler.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ContextHandler(r.Context(), w, r)
	})
}

func enableCorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		next.ServeHTTP(w, r)
	})
}

type homepagePresenter struct {
	Title           string
	GraphiQLEnabled bool
}

func (s *Server) getHomepageHandler() http.HandlerFunc {
	config := s.config
	tmpl, err := template.New("homepage").Parse(`
		<html>
			<h1>{{.Title}}</h1>

			{{if .GraphiQLEnabled}}
				<p>Click here to visit <a href="/graphql">GraphiQL</a>.</p>
			{{end}}
		</html>
	`)

	if err != nil {
		log.Fatalf("Error parsing homepage template err: %s", err)
	}

	return func(w http.ResponseWriter, req *http.Request) {
		presenter := homepagePresenter{
			Title:           "Go Stop Server",
			GraphiQLEnabled: config.GraphiQLEnabled,
		}

		tmpl.Execute(w, presenter)
	}
}
