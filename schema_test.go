package grapher

import (
	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func assertFieldResolver(t *testing.T, fields graphql.Fields, requestString string, expectedValues map[string]interface{}) {
	schemaConfig := graphql.SchemaConfig{}
	schemaConfig.Query = graphql.NewObject(graphql.ObjectConfig{Name: "RootQuery", Fields: fields})
	schema, _ := graphql.NewSchema(schemaConfig)

	assertSchema(t, schema, requestString, expectedValues)
}

func assertSchema(t *testing.T, schema graphql.Schema, requestString string, expectedValues map[string]interface{}) {
	params := graphql.Params{Schema: schema, RequestString: requestString}
	r := graphql.Do(params)

	assert.Empty(t, r.Errors)
	assert.EqualValues(t, expectedValues, r.Data)
}
