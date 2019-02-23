/*
Server implementation of the board game Go.
*/
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/camirmas/go_stop/models"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

type ServerConfig struct {
	host       string
	port       int
	signingKey []byte
}

type Server struct {
	config *ServerConfig
	db     models.DB
	schema *graphql.Schema
}

func NewServer(config *ServerConfig, db models.DB, schema *graphql.Schema) *Server {
	return &Server{
		config,
		db,
		schema,
	}
}

func (s *Server) Start() {
	h := handler.New(&handler.Config{
		Schema:   s.schema,
		Pretty:   true,
		GraphiQL: true,
	})

	log.Println("Listening on http://localhost:8000")
	http.Handle("/graphql", s.applyMiddlewares(h))
	http.HandleFunc("/", Homepage)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func (s *Server) applyMiddlewares(next *handler.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		ctx := authHandler(r.Context(), r)
		next.ContextHandler(ctx, w, r)
	})
}

func authHandler(ctx context.Context, r *http.Request) context.Context {
	auth := r.Header.Get("Authorization")

	if auth != "" {
		s := strings.Split(auth, " ")
		if len(s) == 2 {
			if s[0] == "Bearer" {
				token := s[1]

				ctx = context.WithValue(ctx, "token", token)
			}
		}
	}
	return ctx
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func Homepage(w http.ResponseWriter, req *http.Request) {
	heading := "Go Stop Server"
	greeting := `
		Click here to visit <a href="/graphql">GraphiQL</a>
	`

	fmt.Fprintf(w, "<html><h1>%s</h1><p>%s</p></html>", heading, greeting)
}
