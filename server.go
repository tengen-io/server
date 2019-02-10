/*
Server implementation of the board game Go.
*/
package main

import (
	"context"
	"github.com/camirmas/go_stop/models"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"log"
	"net/http"
	"strings"
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
	http.Handle("/graphql", applyMiddlewares(s, h))
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func applyMiddlewares(s *Server, next *handler.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		ctx := context.WithValue(r.Context(), "db", s.db)
		ctx = context.WithValue(ctx, "signingKey", s.config.signingKey)
		ctx = authHandler(ctx, r)
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
