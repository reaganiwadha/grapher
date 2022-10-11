package grapher

import (
	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func assertSchema(t *testing.T, schema graphql.Schema, requestString string, expectedValues map[string]interface{}) {
	params := graphql.Params{Schema: schema, RequestString: requestString}
	r := graphql.Do(params)

	assert.Empty(t, r.Errors)
	assert.EqualValues(t, expectedValues, r.Data)
}
