package main

import (
	_ "fmt"

	"github.com/camirmas/go_stop/resolvers"
	"github.com/graphql-go/graphql"
)

func CreateSchema() (graphql.Schema, error) {

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
				Resolve: resolvers.GetUser,
			},

			"currentUser": &graphql.Field{
				Type:    userType,
				Resolve: resolvers.CurrentUser,
			},

			"game": &graphql.Field{
				Type: gameType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
				Resolve: resolvers.GetGame,
			},

			"games": &graphql.Field{
				Type: graphql.NewList(gameType),
				Args: graphql.FieldConfigArgument{
					"userId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
				Resolve: resolvers.GetGames,
			},

			"lobby": &graphql.Field{
				Type:    graphql.NewList(gameType),
				Resolve: resolvers.GetLobby,
			},
		},
	})

	schemaConfig := graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	}

	return graphql.NewSchema(schemaConfig)
}
