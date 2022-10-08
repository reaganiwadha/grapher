package grapher

import (
	"github.com/graphql-go/graphql"
)

// BuildRootSchemaConfig Configuration for building the root schema, currently a placeholder
type BuildRootSchemaConfig struct {
}

// MutationQueryCollector A simple Collector struct for *graphql.Field maps for queries and mutations
type MutationQueryCollector interface {
	PushMutations(mutations graphql.Fields)
	PushQueries(queries graphql.Fields)

	BuildRootSchema(optionalCfg ...BuildRootSchemaConfig) (schema graphql.Schema, err error)

	Mutations() graphql.Fields
	Queries() graphql.Fields
}

type mutationQueryCollector struct {
	mutations graphql.Fields
	queries   graphql.Fields
}

// NewMutationQueryCollector Returns a new mutationQueryCollector
func NewMutationQueryCollector() MutationQueryCollector {
	return &mutationQueryCollector{
		mutations: graphql.Fields{},
		queries:   graphql.Fields{},
	}
}

// PushQueries Pushes the queries into the collection
func (m *mutationQueryCollector) PushQueries(queries graphql.Fields) {
	for k, v := range queries {
		m.queries[k] = v
	}
}

// PushMutations Pushes the mutations into the collection
func (m *mutationQueryCollector) PushMutations(mutations graphql.Fields) {
	for k, v := range mutations {
		m.mutations[k] = v
	}
}

// BuildRootSchema Builds the root schema with RootMutation and RootQuery, with the collected mutations and queries
func (m *mutationQueryCollector) BuildRootSchema(optionalCfg ...BuildRootSchemaConfig) (schema graphql.Schema, err error) {
	schemaConfig := graphql.SchemaConfig{}

	if len(m.mutations) != 0 {
		schemaConfig.Mutation = graphql.NewObject(graphql.ObjectConfig{Name: "RootMutation", Fields: m.mutations})
	}

	schemaConfig.Query = graphql.NewObject(graphql.ObjectConfig{Name: "RootQuery", Fields: m.queries})

	return graphql.NewSchema(schemaConfig)
}

func (m *mutationQueryCollector) Mutations() graphql.Fields {
	return m.mutations
}

func (m *mutationQueryCollector) Queries() graphql.Fields {
	return m.queries
}
