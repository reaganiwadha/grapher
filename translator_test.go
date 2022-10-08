package grapher

import (
	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Customer struct {
	ID      int
	Name    string
	Address *string
}

type Order struct {
	ID       int      `json:"id"`
	Customer Customer `json:"customer"`
}

func TestMustTranslate_ValueType(t *testing.T) {
	g := New()

	obj := g.MustTranslateOutput("test")
	assert.Equal(t, graphql.NewNonNull(graphql.String), obj)
}

func TestTranslate_ValueType(t *testing.T) {
	g := New()

	obj, err := g.TranslateOutput("test")
	assert.Equal(t, graphql.NewNonNull(graphql.String), obj)
	assert.NoError(t, err)
}

func TestTranslate_SliceType(t *testing.T) {
	g := New()

	obj, err := g.TranslateOutput([]int{1, 2, 3})
	assert.NoError(t, err)
	assert.Equal(t, graphql.NewList(graphql.NewNonNull(graphql.Int)), obj)
}

func TestTranslate_ArrayType(t *testing.T) {
	g := New()

	obj, err := g.TranslateOutput([3]int{3, 3, 2})
	assert.NoError(t, err)
	assert.Equal(t, graphql.NewList(graphql.NewNonNull(graphql.Int)), obj)
}

func TestTranslate_MapType(t *testing.T) {
	g := New()

	obj, err := g.TranslateOutput(map[string]string{
		"cool": "hi",
	})
	assert.NoError(t, err)
	assert.Equal(t, graphql.NewNonNull(ScalarJSON), obj)
}

func testCustomerType(t *testing.T, actual graphql.Output) {
	obj, ok := actual.(*graphql.Object)
	if !ok {
		obj = actual.(*graphql.NonNull).OfType.(*graphql.Object)
	}

	assert.Equal(t, "Customer", obj.Name())
	assert.Equal(t, obj.Fields()["ID"].Type, graphql.NewNonNull(graphql.Int))
	assert.Equal(t, obj.Fields()["Name"].Type, graphql.NewNonNull(graphql.String))
	assert.Equal(t, obj.Fields()["Address"].Type, graphql.String)
}

func testOrderType(t *testing.T, actual graphql.Output) {
	obj := actual.(*graphql.Object)
	assert.Equal(t, "Order", obj.Name())
	assert.Equal(t, obj.Fields()["id"].Type, graphql.NewNonNull(graphql.Int))
	testCustomerType(t, obj.Fields()["customer"].Type)
}

func TestTranslate_SimpleStruct(t *testing.T) {
	g := New()

	obj, err := g.TranslateOutput(&Customer{})
	assert.NoError(t, err)
	testCustomerType(t, obj)
}

func TestTranslate_NestedStruct_WithoutCache(t *testing.T) {
	g := New()

	obj, err := g.TranslateOutput(&Order{})
	assert.NoError(t, err)
	testOrderType(t, obj)
}

func TestTranslate_NestedStruct_WithCache(t *testing.T) {
	g := New()

	obj, err := g.TranslateOutput(&Customer{})
	assert.NoError(t, err)
	testCustomerType(t, obj)

	orderObj, err := g.TranslateOutput(&Order{})

	assert.EqualValues(t, "Customer!", orderObj.(*graphql.Object).Fields()["customer"].Type.Name())
	assert.NoError(t, err)
	testOrderType(t, orderObj)
}
