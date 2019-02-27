package main

import (
	_ "fmt"
	"github.com/graphql-go/graphql"
	"github.com/tengen-io/server/resolvers"
)

func NewSchema(resolvers *resolvers.Resolvers) (graphql.Schema, error) {
	objects := buildObjects()
	schemaConfig := graphql.SchemaConfig{
		Query:    buildQueries(resolvers, objects),
		Mutation: buildMutations(resolvers, objects),
	}

	return graphql.NewSchema(schemaConfig)
}
