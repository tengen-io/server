package gql

import (
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type ContextKey = byte

const (
	IdentityContextKey ContextKey = iota
)

func (s *server) LoginHandler() http.Handler {
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
		identity, err := s.auth.checkPasswordByEmail(credentials.Email, credentials.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		token, err := s.auth.signJWT(*identity)

		w.Header().Set("Content-Type", "application/json")
		out, err := json.Marshal(struct{ Token string }{token})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(out)
	})
}

func (s *server) VerifyTokenMiddleware(next http.Handler) http.Handler {
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
			token, err := s.auth.validateJWT(tokenStr)
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

			identity, err := s.repo.GetIdentityById(int32(idInt))
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

func (s *server) RegistrationHandler() http.Handler {
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
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(in.Password), s.auth.bcryptCost)
		identity, err := s.repo.CreateIdentity(in.Email, passwordHash, in.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token, err := s.auth.signJWT(*identity)
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
