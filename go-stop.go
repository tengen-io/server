/*
Server implementation of the board game Go.
*/
package main

import (
	"context"
	"database/sql"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"log"
	"net/http"
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
	http.Handle("/graphql", DBMiddleware(db, h))
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func DBMiddleware(db *sql.DB, next *handler.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "db", db)

		next.ContextHandler(ctx, w, r)
	})
}

func main() {
	s := Server{}
	s.Start()
}
