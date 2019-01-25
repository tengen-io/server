package main

import (
	"github.com/camirmas/go_stop/resolvers"
	"github.com/graphql-go/graphql"
)

func buildMutations(objects *Objects) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createUser": &graphql.Field{
				Type: objects.AuthUser,
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
				Resolve: resolvers.CreateUser,
			},

			"createGame": &graphql.Field{
				Type: objects.Game,
				Args: graphql.FieldConfigArgument{
					"opponentUsername": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: resolvers.CreateGame,
			},

			"pass": &graphql.Field{
				Type: objects.Game,
				Args: graphql.FieldConfigArgument{
					"gameId": &graphql.ArgumentConfig{
						Type: graphql.ID,
					},
				},
				Resolve: resolvers.Pass,
			},

			"logIn": &graphql.Field{
				Type: objects.AuthUser,
				Args: graphql.FieldConfigArgument{
					"username": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: resolvers.LogIn,
			},

			"addStone": &graphql.Field{
				Type: objects.Game,
				Args: graphql.FieldConfigArgument{
					"gameId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
					"x": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"y": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: resolvers.AddStone,
			},
		},
	})
}
