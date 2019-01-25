package main

import (
	"github.com/camirmas/go_stop/resolvers"
	"github.com/graphql-go/graphql"
)

func buildQueries(objects *Objects) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"user": &graphql.Field{
				Type: objects.User,
				Args: graphql.FieldConfigArgument{
					"username": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: resolvers.GetUser,
			},

			"currentUser": &graphql.Field{
				Type:    objects.User,
				Resolve: resolvers.CurrentUser,
			},

			"game": &graphql.Field{
				Type: objects.Game,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
				Resolve: resolvers.GetGame,
			},

			"games": &graphql.Field{
				Type: graphql.NewList(objects.Game),
				Args: graphql.FieldConfigArgument{
					"userId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
				Resolve: resolvers.GetGames,
			},

			"lobby": &graphql.Field{
				Type:    graphql.NewList(objects.Game),
				Resolve: resolvers.GetLobby,
			},
		},
	})
}
