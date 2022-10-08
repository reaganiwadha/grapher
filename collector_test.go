package grapher

import (
	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testMutationCollection1 = graphql.Fields{
	"mutation1": {
		Type: graphql.String,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return "mutation1", nil
		},
	},
}

var testMutationCollection2 = graphql.Fields{
	"mutation2": {
		Type: graphql.String,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return "mutation2", nil
		},
	},
}

var testQueryCollection1 = graphql.Fields{
	"query1": {
		Type: graphql.String,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return "query1", nil
		},
	},
}
var testQueryCollection2 = graphql.Fields{
	"query2": {
		Type: graphql.String,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return "query2", nil
		},
	},
}

func TestMutationQueryCollector_ProperlyInitialized(t *testing.T) {
	mqc := NewMutationQueryCollector()

	assert.NotNil(t, mqc.Queries())
	assert.NotNil(t, mqc.Mutations())
}

func TestMutationQueryCollector_QueryOnly(t *testing.T) {
	mqc := NewMutationQueryCollector()

	mqc.PushQueries(testQueryCollection1)
	mqc.PushQueries(testQueryCollection2)
	schema, err := mqc.BuildRootSchema()

	assert.NoError(t, err)
	assert.NotNil(t, schema)

	assertSchema(t, schema, `
		query {
			query1
			query2
		}
	`, map[string]interface{}{
		"query1": "query1",
		"query2": "query2",
	})
}

func TestMutationQueryCollector_WithMutation(t *testing.T) {
	mqc := NewMutationQueryCollector()

	mqc.PushMutations(testMutationCollection1)
	mqc.PushMutations(testMutationCollection2)
	mqc.PushQueries(testQueryCollection1)
	mqc.PushQueries(testQueryCollection2)

	schema, err := mqc.BuildRootSchema()

	assert.NoError(t, err)
	assert.NotNil(t, schema)

	assertSchema(t, schema, `
		mutation {
			mutation1
			mutation2
		}
	`, map[string]interface{}{
		"mutation1": "mutation1",
		"mutation2": "mutation2",
	})

	assertSchema(t, schema, `
		query {
			query1
			query2
		}
	`, map[string]interface{}{
		"query1": "query1",
		"query2": "query2",
	})
}

func assertSchema(t *testing.T, schema graphql.Schema, requestString string, expectedValues map[string]interface{}) {
	params := graphql.Params{Schema: schema, RequestString: requestString}
	r := graphql.Do(params)

	assert.Empty(t, r.Errors)
	assert.EqualValues(t, expectedValues, r.Data)
}
