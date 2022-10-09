package scalars

import (
	"github.com/graphql-go/graphql/language/ast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJSONScalar_ParseValue(t *testing.T) {
	assert.EqualValues(t, "test", ScalarJSON.ParseValue("test"))
}

func TestJSONScalar_Serialize(t *testing.T) {
	assert.EqualValues(t, "test", ScalarJSON.Serialize("test"))
}

func TestJSONScalar_ParseLiteral_String(t *testing.T) {
	astVal := ast.NewStringValue(&ast.StringValue{
		Value: "test",
	})

	assert.EqualValues(t, "test", ScalarJSON.ParseLiteral(astVal))
}

func TestJSONScalar_ParseLiteral_Object(t *testing.T) {
	val := map[string]interface{}{
		"test": "test",
	}

	astVal := ast.NewObjectValue(&ast.ObjectValue{
		Fields: []*ast.ObjectField{
			{
				Name: &ast.Name{
					Value: "test",
				},
				Value: ast.NewStringValue(&ast.StringValue{
					Value: "test",
				}),
			},
		},
	})

	assert.EqualValues(t, val, ScalarJSON.ParseLiteral(astVal))
}

func TestJSONScalar_ParseLiteral_List(t *testing.T) {
	val := []interface{}{
		"test1",
		"test2",
	}

	astVal := ast.NewListValue(&ast.ListValue{
		Values: []ast.Value{
			ast.NewStringValue(&ast.StringValue{
				Value: "test1",
			}),
			ast.NewStringValue(&ast.StringValue{
				Value: "test2",
			}),
		},
	})

	assert.EqualValues(t, val, ScalarJSON.ParseLiteral(astVal))
}

func TestJSONScalar_ParseLiteral_NilOnUnkown(t *testing.T) {
	assert.EqualValues(t, nil, ScalarJSON.ParseLiteral(ast.NewEnumValue(&ast.EnumValue{})))
}
