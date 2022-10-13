package grapher

import (
	"github.com/go-playground/validator/v10"
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

type CustomArgParserStruct struct {
	Counter int `json:"counter" validator:"max=10"`
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

func TestFieldBuilder_MustBuild_Resolver(t *testing.T) {
	ret := NewFieldBuilder[NoArgs, string]().
		WithDescription("CoolTest").
		WithResolver(func(p graphql.ResolveParams, args NoArgs) (string, error) {
			return "hi", nil
		}).
		MustBuild()

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

func TestFieldBuilder_MustBuild_InvalidArgs(t *testing.T) {
	assert.Panics(t, func() {
		NewFieldBuilder[string, string]().
			WithResolver(func(p graphql.ResolveParams, args string) (string, error) {
				return "", nil
			}).
			MustBuild()
	})
}

func TestFieldBuilder_Middlewares(t *testing.T) {
	middleware1 := func(nextFn graphql.FieldResolveFn) graphql.FieldResolveFn {
		return func(p graphql.ResolveParams) (interface{}, error) {
			return nextFn(p)
		}
	}

	middleware2 := func(nextFn graphql.FieldResolveFn) graphql.FieldResolveFn {
		return func(p graphql.ResolveParams) (interface{}, error) {
			return nextFn(p)
		}
	}

	ret, err := NewFieldBuilder[NoArgs, string]().
		AddMiddleware(middleware1).
		AddMiddleware(middleware2).
		WithResolver(func(p graphql.ResolveParams, arg NoArgs) (string, error) {
			return "ok", nil
		}).Build()

	assert.NoError(t, err)

	fields := graphql.Fields{
		"query1": ret,
	}

	assertFieldResolver(t, fields, `
		query {
			query1
		}
	`, map[string]interface{}{
		"query1": "ok",
	})
}

func TestFieldBuilder_CustomArgParser(t *testing.T) {
	v := validator.New()

	argParser := func(arg any) error {
		return v.Struct(&arg)
	}

	ret, err := NewFieldBuilder[CustomArgParserStruct, string]().
		WithCustomArgValidator(argParser).
		WithResolver(func(p graphql.ResolveParams, arg CustomArgParserStruct) (string, error) {
			assert.Fail(t, "resolver called")
			return "ok", nil
		}).Build()

	assert.NoError(t, err)

	fields := graphql.Fields{
		"query1": ret,
	}

	schemaConfig := graphql.SchemaConfig{}
	schemaConfig.Query = graphql.NewObject(graphql.ObjectConfig{Name: "RootQuery", Fields: fields})
	schema, _ := graphql.NewSchema(schemaConfig)

	params := graphql.Params{Schema: schema, RequestString: `query { query1(counter : 44) }`}
	r := graphql.Do(params)

	assert.NotEmpty(t, r.Errors)
}
