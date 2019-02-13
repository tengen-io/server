package main

import (
	_ "fmt"
	"github.com/camirmas/go_stop/resolvers"
	"github.com/graphql-go/graphql"
)

func NewSchema(resolvers *resolvers.Resolvers) (graphql.Schema, error) {
	objects := buildObjects()
	schemaConfig := graphql.SchemaConfig{
		Query:    buildQueries(resolvers, objects),
		Mutation: buildMutations(resolvers, objects),
	}

	return graphql.NewSchema(schemaConfig)
}
