/*
Server implementation of the board game Go.
*/
package main

import (
	"context"
	"encoding/base64"
	_ "fmt"
	"github.com/camirmas/go_stop/models"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Server struct {
	db         models.Database
	schema     *graphql.Schema
	signingKey []byte
}

func (s *Server) Start() {
	db, err := models.ConnectDB()
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
	s.signingKey = getKey()

	h := handler.New(&handler.Config{
		Schema:   &schema,
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
		ctx = context.WithValue(ctx, "signingKey", s.signingKey)
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

func getKey() []byte {
	filename := os.Getenv("JWT_SECRET_KEY")
	secret := getSecret(filename)
	encoded, _ := base64.StdEncoding.DecodeString(secret)

	return encoded
}

func getSecret(fileName string) string {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return ""
	}

	return string(file)
}

func main() {
	s := Server{}
	s.Start()
}
