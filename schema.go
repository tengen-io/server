package main

import (
	"database/sql"
	"github.com/graphql-go/graphql"
)

func CreateSchema() (graphql.Schema, error) {
	userType := graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if user, ok := p.Source.(*User); ok {
						return user.Id, nil
					}
					return nil, nil
				},
			},
			"username": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if user, ok := p.Source.(*User); ok {
						return user.Username, nil
					}
					return nil, nil
				},
			},
			"email": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if user, ok := p.Source.(*User); ok {
						return user.Email, nil
					}
					return nil, nil
				},
			},
		},
	})

	tokenType := graphql.NewObject(graphql.ObjectConfig{
		Name: "AuthUser",
		Fields: graphql.Fields{
			"user": &graphql.Field{
				Type: userType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if authUser, ok := p.Source.(*AuthUser); ok {
						return authUser.User, nil
					}
					return nil, nil
				},
			},
			"token": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if authUser, ok := p.Source.(*AuthUser); ok {
						return authUser.Jwt, nil
					}
					return nil, nil
				},
			},
		},
	})

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"user": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"username": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					db := p.Context.Value("db").(*sql.DB)
					return GetUser(db, p.Args["username"].(string))
				},
			},

			"currentUser": &graphql.Field{
				Type: userType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					db := p.Context.Value("db").(*sql.DB)

					token, ok := p.Context.Value("token").(string)

					if !ok {
						return nil, missingTokenError{}
					}

					claims, err := ValidateToken(token)

					if err != nil {
						return nil, err
					}

					return GetUser(db, claims.Username)
				},
			},
		},
	})

	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createUser": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"username": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"passwordConfirmation": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					db := p.Context.Value("db").(*sql.DB)
					username := p.Args["username"].(string)
					email := p.Args["email"].(string)
					password := p.Args["password"].(string)
					passwordConfirm := p.Args["passwordConfirmation"].(string)
					return CreateUser(db, username, email, password, passwordConfirm)
				},
			},

			"logIn": &graphql.Field{
				Type: tokenType,
				Args: graphql.FieldConfigArgument{
					"username": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					username := p.Args["username"].(string)
					password := p.Args["password"].(string)
					db := p.Context.Value("db").(*sql.DB)
					user, err := CheckPw(db, username, password)

					if err != nil {
						return nil, err
					}

					token, err := GenerateToken(username)

					if err != nil {
						return user, err
					}

					authUser := &AuthUser{token, user}

					return authUser, nil
				},
			},
		},
	})

	schemaConfig := graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	}

	return graphql.NewSchema(schemaConfig)
}

type missingTokenError struct{}

func (e missingTokenError) Error() string {
	return "Missing Authorization header"
}
