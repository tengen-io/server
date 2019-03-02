package server

import (
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func (s *Server) LoginHandler() http.Handler {
	type credentials struct {
		Email string
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
			}

			claims, ok := token.Claims.(*jwt.StandardClaims)

			if !ok {
				http.Error(w, "unable to cast claims", http.StatusInternalServerError)
				return
			}

			idInt, err := strconv.Atoi(claims.Id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			identity, err := s.identity.GetIdentityById(int32(idInt))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			ctx = context.WithValue(ctx, "identity", identity)
		}

		next.ServeHTTP(w, r)
	})
}
