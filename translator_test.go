package grapher

import (
	"github.com/graphql-go/graphql"
	"github.com/reaganiwadha/grapher/scalars"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMustTranslate_ValueType(t *testing.T) {
	g := New()

	obj := g.MustTranslate("test")
	assert.Equal(t, graphql.NewNonNull(graphql.String), obj)
}

func TestTranslate_ValueType(t *testing.T) {
	g := New()

	obj, err := g.Translate("test")
	assert.Equal(t, graphql.NewNonNull(graphql.String), obj)
	assert.NoError(t, err)
}

func TestTranslate_SliceType(t *testing.T) {
	g := New()

	obj, err := g.Translate([]int{1, 2, 3})
	assert.NoError(t, err)
	assert.Equal(t, graphql.NewList(graphql.NewNonNull(graphql.Int)), obj)
}

func TestTranslate_ArrayType(t *testing.T) {
	g := New()

	obj, err := g.Translate([3]int{3, 3, 2})
	assert.NoError(t, err)
	assert.Equal(t, graphql.NewList(graphql.NewNonNull(graphql.Int)), obj)
}

func TestTranslate_MapType(t *testing.T) {
	g := New()

	obj, err := g.Translate(map[string]string{
		"cool": "hi",
	})
	assert.NoError(t, err)
	assert.Equal(t, graphql.NewNonNull(scalars.ScalarJSON), obj)
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

type Customer struct {
	ID      int
	Name    string
	Address *string
}

type Order struct {
	ID       int      `json:"id"`
	Customer Customer `json:"customer"`
}

func TestTranslate_SimpleStruct(t *testing.T) {
	g := New()

	obj, err := g.Translate(&Customer{})
	assert.NoError(t, err)
	testCustomerType(t, obj)
}

func TestTranslate_NestedStruct_WithoutCache(t *testing.T) {
	g := New()

	obj, err := g.Translate(&Order{})
	assert.NoError(t, err)
	testOrderType(t, obj)
}

func TestTranslate_NestedStruct_WithCache(t *testing.T) {
	g := New()

	obj, err := g.Translate(&Customer{})
	assert.NoError(t, err)
	testCustomerType(t, obj)

	orderObj, err := g.Translate(&Order{})

	assert.EqualValues(t, "Customer!", orderObj.(*graphql.Object).Fields()["customer"].Type.Name())
	assert.NoError(t, err)
	testOrderType(t, orderObj)
}

type QueryUserWhereClause struct {
	NameLike string `json:"name_like"`
}

type QueryUser struct {
	Limit  *int                  `json:"limit"`
	Offset *int                  `json:"offset"`
	Where  *QueryUserWhereClause `json:"where"`
}

func testQueryUserWhereClause(t *testing.T, obj *graphql.InputObject) {
	fields := obj.Fields()
	assert.Equal(t, fields["name_like"].Type, graphql.NewNonNull(graphql.String))
}

func testQueryUser(t *testing.T, obj *graphql.InputObject) {
	fields := obj.Fields()
	assert.Equal(t, "QueryUser", obj.Name())
	assert.Equal(t, fields["limit"].Type, graphql.Int)
	assert.Equal(t, fields["offset"].Type, graphql.Int)

	testQueryUserWhereClause(t, fields["where"].Type.(*graphql.InputObject))
}

func TestTranslateInput_NestedStruct(t *testing.T) {
	g := New()

	inputObj, err := g.TranslateInputObject(&QueryUser{})
	assert.NoError(t, err)
	testQueryUser(t, inputObj)
}

func TestTranslateInput_NestedStruct_WithCache(t *testing.T) {
	g := New()

	whereObj, err := g.TranslateInputObject(&QueryUserWhereClause{})
	testQueryUserWhereClause(t, whereObj)

	queryObj, err := g.TranslateInputObject(&QueryUser{})
	assert.NoError(t, err)
	testQueryUser(t, queryObj)
}

func TestTranslateInput_Value_Errors(t *testing.T) {
	g := New()

	_, err := g.TranslateInputObject(2)

	assert.Error(t, err)
}
