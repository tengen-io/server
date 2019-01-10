package main

import (
	"database/sql"
	_ "fmt"
	"github.com/graphql-go/graphql"
)

func CreateSchema() (graphql.Schema, error) {
	userType := graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
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

	playerType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Player",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if player, ok := p.Source.(Player); ok {
						return player.Id, nil
					}
					return nil, nil
				},
			},
			"user": &graphql.Field{
				Type: userType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if player, ok := p.Source.(Player); ok {
						return player.User, nil
					}
					return nil, nil
				},
			},
		},
	})

	gameType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Game",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if game, ok := p.Source.(*Game); ok {
						return game.Id, nil
					}
					return nil, nil
				},
			},
			"status": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if game, ok := p.Source.(*Game); ok {
						return game.Status, nil
					}
					return nil, nil
				},
			},
			"playerTurnId": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if game, ok := p.Source.(*Game); ok {
						return game.PlayerTurnId, nil
					}
					return nil, nil
				},
			},
			"players": &graphql.Field{
				Type: graphql.NewList(playerType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if game, ok := p.Source.(*Game); ok {
						return game.Players, nil
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

					return GetUser(db, claims.UserId)
				},
			},

			"game": &graphql.Field{
				Type: gameType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					db := p.Context.Value("db").(*sql.DB)
					return GetGame(db, p.Args["id"].(int))
				},
			},

			"games": &graphql.Field{
				Type: graphql.NewList(gameType),
				Args: graphql.FieldConfigArgument{
					"userId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					db := p.Context.Value("db").(*sql.DB)
					return GetGames(db, p.Args["userId"].(int))
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

			"createGame": &graphql.Field{
				Type: gameType,
				Args: graphql.FieldConfigArgument{
					"opponentId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					db := p.Context.Value("db").(*sql.DB)
					token := p.Context.Value("token").(string)
					opponentId := p.Args["opponentId"].(int)

					claims, err := ValidateToken(token)

					if err != nil {
						return nil, err
					}

					return CreateGame(db, claims.UserId, opponentId)
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

					token, err := GenerateToken(user.Id)

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
