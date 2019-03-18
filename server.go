package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type ContextKey = byte

const (
	IdentityContextKey ContextKey = iota
)

type ServerConfig struct {
	Host            string
	Port            int
	GraphiQLEnabled bool
}

type Server struct {
	config           *ServerConfig
	executableSchema graphql.ExecutableSchema
	auth             *AuthRepository
	identity         *IdentityRepository
}

func NewServer(config *ServerConfig, schema graphql.ExecutableSchema, auth *AuthRepository, identity *IdentityRepository) *Server {
	return &Server{
		config,
		schema,
		auth,
		identity,
	}
}

func (s *Server) Start() {
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

func (s *Server) LoginHandler() http.Handler {
	type credentials struct {
		Email    string
		Password string
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var credentials credentials
		err = json.Unmarshal(b, &credentials)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// TODO(eac): find a way to differentiate between auth and db failure
		identity, err := s.auth.CheckPasswordByEmail(credentials.Email, credentials.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		token, err := s.auth.SignJWT(*identity)

		w.Header().Set("Content-Type", "application/json")
		out, err := json.Marshal(struct{ Token string }{token})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(out)
	})
}

func (s *Server) VerifyTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authHeader := r.Header.Get("Authorization")

		if authHeader != "" {
			authParts := strings.Split(authHeader, " ")
			if len(authParts) != 2 || authParts[0] != "Bearer" {
				http.Error(w, "invalid auth header", http.StatusBadRequest)
				return
			}

			tokenStr := authParts[1]
			token, err := s.auth.ValidateJWT(tokenStr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(*jwt.StandardClaims)

			if !ok {
				http.Error(w, "unable to cast claims", http.StatusInternalServerError)
				return
			}

			idInt, err := strconv.Atoi(claims.Id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			identity, err := s.identity.GetIdentityById(int32(idInt))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			ctx = context.WithValue(ctx, IdentityContextKey, *identity)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) RegistrationHandler() http.Handler {
	type input struct {
		Email    string
		Password string
		Name     string
	}

	type output struct {
		Id    string `json:"id"`
		Token string `json:"token"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var in input
		err = json.Unmarshal(b, &in)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if in.Email == "" || in.Name == "" || in.Password == "" {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		// TODO(eac): Distinguish between database errors?
		identity, err := s.identity.CreateIdentity(in.Email, in.Password, in.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token, err := s.auth.SignJWT(*identity)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		out, err := json.Marshal(output{Id: identity.Id, Token: token})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(out)
	})
}
