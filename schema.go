package main

import (
	_ "fmt"
	"github.com/graphql-go/graphql"
)

func CreateSchema() (graphql.Schema, error) {
	objects := buildObjects()
	schemaConfig := graphql.SchemaConfig{
		Query:    buildQueries(objects),
		Mutation: buildMutations(objects),
	}

	return graphql.NewSchema(schemaConfig)
}
