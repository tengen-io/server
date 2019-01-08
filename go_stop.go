/*
Server implementation of the board game Go.
*/
package main

import (
	"context"
	"database/sql"
	_ "fmt"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"log"
	"net/http"
	"strings"
)

type Server struct {
	db     *sql.DB
	schema *graphql.Schema
}

func (s *Server) Start() {
	db, err := ConnectDB()
	if err != nil {
		log.Fatal(err)
		return
	} else {
		log.Println("Connected to postgres")
	}
	s.db = db

	schema, err := CreateSchema()
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	} else {
		log.Println("Created GraphQL Schema")
	}
	s.schema = &schema

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	log.Println("Listening on http://localhost:8000")
	http.Handle("/graphql", applyMiddlewares(db, h))
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func applyMiddlewares(db *sql.DB, next *handler.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "db", db)
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

func main() {
	s := Server{}
	s.Start()
}
