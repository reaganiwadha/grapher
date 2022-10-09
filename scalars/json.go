package scalars

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
)

// ScalarJSON Taken from github.com/bhoriuchi/graphql-go-tools/blob/master/scalars/json.go
var ScalarJSON = graphql.NewScalar(graphql.ScalarConfig{
	Name: "JSON",
	Serialize: func(value interface{}) interface{} {
		return value
	},
	ParseValue: func(value interface{}) interface{} {
		return value
	},
	ParseLiteral: parseLiteralJSON,
})

func parseLiteralJSON(astValue ast.Value) interface{} {
	switch kind := astValue.GetKind(); kind {
	case kinds.StringValue, kinds.BooleanValue, kinds.IntValue, kinds.FloatValue:
		return astValue.GetValue()

	case kinds.ObjectValue:
		obj := make(map[string]interface{})
		for _, v := range astValue.GetValue().([]*ast.ObjectField) {
			obj[v.Name.Value] = parseLiteralJSON(v.Value)
		}
		return obj

	case kinds.ListValue:
		list := make([]interface{}, 0)
		for _, v := range astValue.GetValue().([]ast.Value) {
			list = append(list, parseLiteralJSON(v))
		}
		return list

	default:
		return nil
	}
}
