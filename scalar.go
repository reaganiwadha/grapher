package grapher

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

var ScalarJSON = graphql.NewScalar(graphql.ScalarConfig{
	Name: "JSON",
	Serialize: func(value interface{}) interface{} {
		return value
	},
	ParseValue: func(value interface{}) interface{} {
		return value
	},
	ParseLiteral: func(valueAST ast.Value) interface{} {
		// TODO
		return nil
	},
})
