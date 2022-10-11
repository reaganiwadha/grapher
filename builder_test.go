package grapher

import (
	"github.com/graphql-go/graphql"
	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
	"testing"
)

type ArgStruct struct {
	Name string `json:"name"`
}

type CustomTranslatorStruct struct {
	Name    null.String `json:"name"`
	Counter null.Int    `json:"counter"`
}

func TestFieldBuilder_Resolver(t *testing.T) {
	ret, err := NewFieldBuilder[NoArgs, string]().
		WithDescription("CoolTest").
		WithResolver(func(p graphql.ResolveParams, args NoArgs) (string, error) {
			return "hi", nil
		}).
		Build()

	assert.NoError(t, err)
	assert.NotNil(t, ret)

	fields := graphql.Fields{
		"query1": ret,
	}

	assert.Equal(t, "CoolTest", ret.Description)
	assertFieldResolver(t, fields, `
		query {
			query1
		}
	`, map[string]interface{}{
		"query1": "hi",
	})
}

func TestFieldBuilder_CustomTranslator(t *testing.T) {
	translator := NewTranslator(&TranslatorConfig{
		PredefinedTranslation: TranslationMap{
			"null.String": graphql.String,
			"null.Int":    graphql.Int,
		},
	})

	ret, err := NewFieldBuilder[CustomTranslatorStruct, string](FieldBuilderConfig{Translator: translator}).
		WithDescription("CoolTest").
		WithResolver(func(p graphql.ResolveParams, args CustomTranslatorStruct) (string, error) {
			assert.False(t, args.Name.Valid)
			assert.Equal(t, int64(4), args.Counter.Int64)
			return "hi", nil
		}).
		Build()

	assert.NoError(t, err)
	assert.NotNil(t, ret)

	fields := graphql.Fields{
		"query1": ret,
	}

	assert.Equal(t, "CoolTest", ret.Description)
	assertFieldResolver(t, fields, `
		query {
			query1(counter : 4)
		}
	`, map[string]interface{}{
		"query1": "hi",
	})
}

func TestFieldBuilder_WithArgs(t *testing.T) {
	ret, err := NewFieldBuilder[ArgStruct, string]().
		WithDescription("CoolTest").
		WithResolver(func(p graphql.ResolveParams, args ArgStruct) (string, error) {
			assert.Equal(t, args.Name, "jake")
			return args.Name, nil
		}).
		Build()

	assert.NoError(t, err)
	assert.NotNil(t, ret)

	fields := graphql.Fields{
		"query1": ret,
	}

	assert.Equal(t, "CoolTest", ret.Description)
	assertFieldResolver(t, fields, `
		query {
			query1(name : "jake")
		}
	`, map[string]interface{}{
		"query1": "jake",
	})
}

func TestFieldBuilder_InvalidArgs(t *testing.T) {
	_, err := NewFieldBuilder[string, string]().
		WithResolver(func(p graphql.ResolveParams, args string) (string, error) {
			return "", nil
		}).
		Build()

	assert.Error(t, err)
}

func assertFieldResolver(t *testing.T, fields graphql.Fields, requestString string, expectedValues map[string]interface{}) {
	schemaConfig := graphql.SchemaConfig{}
	schemaConfig.Query = graphql.NewObject(graphql.ObjectConfig{Name: "RootQuery", Fields: fields})
	schema, _ := graphql.NewSchema(schemaConfig)

	assertSchema(t, schema, requestString, expectedValues)
}
