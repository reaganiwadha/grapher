package grapher

import (
	"github.com/graphql-go/graphql"
	"github.com/guregu/null"
	"github.com/reaganiwadha/grapher/scalars"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTranslate_ValueType(t *testing.T) {
	g := NewTranslator()

	obj, err := g.Translate("test")
	assert.Equal(t, graphql.NewNonNull(graphql.String), obj)
	assert.NoError(t, err)
}

func TestTranslate_ValueType_WithPredefinedTranslation(t *testing.T) {
	g := NewTranslator(&TranslatorConfig{
		PredefinedTranslation: TranslationMap{
			"null.Int": graphql.Int,
		},
	})

	obj, err := g.Translate(null.IntFrom(1))
	assert.NoError(t, err)
	assert.Equal(t, graphql.Int, obj)
}

func TestTranslate_SliceType(t *testing.T) {
	g := NewTranslator()

	obj, err := g.Translate([]int{1, 2, 3})
	assert.NoError(t, err)
	assert.Equal(t, graphql.NewList(graphql.NewNonNull(graphql.Int)), obj)
}

func TestTranslate_ArrayType(t *testing.T) {
	g := NewTranslator()

	obj, err := g.Translate([3]int{3, 3, 2})
	assert.NoError(t, err)
	assert.Equal(t, graphql.NewList(graphql.NewNonNull(graphql.Int)), obj)
}

func TestTranslate_MapType(t *testing.T) {
	g := NewTranslator()

	obj, err := g.Translate(map[string]string{
		"cool": "hi",
	})
	assert.NoError(t, err)
	assert.Equal(t, graphql.NewNonNull(scalars.ScalarJSON), obj)
}

func TestMustTranslate_ValueType(t *testing.T) {
	g := NewTranslator()

	obj := g.MustTranslate("test")
	assert.Equal(t, graphql.NewNonNull(graphql.String), obj)
}

func testCustomerType(t *testing.T, actual graphql.Output) {
	obj, ok := actual.(*graphql.Object)
	if !ok {
		obj = actual.(*graphql.NonNull).OfType.(*graphql.Object)
	}

	assert.Equal(t, "Customer", obj.Name())
	assert.Equal(t, obj.Fields()["ID"].Type, graphql.NewNonNull(graphql.Int))
	assert.Equal(t, obj.Fields()["Name"].Type, graphql.NewNonNull(graphql.String))
	assert.Equal(t, obj.Fields()["Name"].Description, "The name of the customer")
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
	Name    string `grapher_d:"The name of the customer"`
	Address *string
}

type Order struct {
	ID       int      `json:"id"`
	Customer Customer `json:"customer"`
}

func TestTranslator_Translate_SimpleStruct(t *testing.T) {
	g := NewTranslator()

	obj, err := g.Translate(&Customer{})
	assert.NoError(t, err)
	testCustomerType(t, obj)
}

type Generic[T1 any, T2 any] struct {
	ID int `json:"id"`
	T1 T1  `json:"t1"`
	T2 T2  `json:"t2"`
}

func TestTranslator_Translate_GenericStruct(t *testing.T) {
	g := NewTranslator()

	obj, err := g.Translate(&Generic[Customer, Order]{})
	assert.NoError(t, err)
	assert.NoError(t, obj.Error())
	assert.NotNil(t, obj)
}

func TestTranslator_Translate_WithoutCache(t *testing.T) {
	g := NewTranslator()

	obj, err := g.Translate(&Order{})
	assert.NoError(t, err)
	testOrderType(t, obj)
}

func TestTranslator_Translate_NestedStruct_WithCache(t *testing.T) {
	g := NewTranslator()

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

func TestTranslator_TranslateInput_NestedStruct(t *testing.T) {
	g := NewTranslator()

	inputObj, err := g.TranslateInput(&QueryUser{})
	assert.NoError(t, err)
	testQueryUser(t, inputObj)
}

func TestTranslator_TranslateInput_NestedStruct_WithCache(t *testing.T) {
	g := NewTranslator()

	whereObj, err := g.TranslateInput(&QueryUserWhereClause{})
	testQueryUserWhereClause(t, whereObj)

	queryObj, err := g.TranslateInput(&QueryUser{})
	assert.NoError(t, err)
	testQueryUser(t, queryObj)
}

func TestTranslator_TranslateInput_ValueType_Errors(t *testing.T) {
	g := NewTranslator()

	_, err := g.TranslateInput(2)

	assert.Error(t, err)
}

func TestTranslator_MustTranslateInput_ValueType_Panics(t *testing.T) {
	g := NewTranslator()

	assert.Panics(t, func() {
		g.MustTranslateInput(2)
	})
}

func TestTranslator_MustTranslateInput_NestedStruct(t *testing.T) {
	g := NewTranslator()

	obj := g.MustTranslateInput(&QueryUser{})
	testQueryUser(t, obj)
}

func testQueryUserFieldConfigArgument(t *testing.T, fieldArgs graphql.FieldConfigArgument) {
	assert.Equal(t, fieldArgs["limit"].Type, graphql.Int)
	assert.Equal(t, fieldArgs["offset"].Type, graphql.Int)

	whereObj := fieldArgs["where"].Type.(*graphql.InputObject)
	whereArgs := whereObj.Fields()
	assert.Equal(t, whereArgs["name_like"].Type, graphql.NewNonNull(graphql.String))
}

func TestTranslator_TranslateArgs_NestedStruct(t *testing.T) {
	g := NewTranslator()

	fieldArgs, err := g.TranslateArgs(&QueryUser{})
	assert.NoError(t, err)
	testQueryUserFieldConfigArgument(t, fieldArgs)
}

func TestTranslator_TranslateArgs_ValueType_Errors(t *testing.T) {
	g := NewTranslator()

	_, err := g.TranslateArgs(1)
	assert.Error(t, err)
}

func TestTranslator_MustTranslateArgs_NestedStruct(t *testing.T) {
	g := NewTranslator()

	fieldArgs := g.MustTranslateArgs(&QueryUser{})
	testQueryUserFieldConfigArgument(t, fieldArgs)
}

type WithTime struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func TestTranslator_Time(t *testing.T) {
	g := NewTranslator()

	fieldArgs := g.MustTranslateArgs(&WithTime{})

	assert.Equal(t, fieldArgs["created_at"].Type, graphql.NewNonNull(graphql.DateTime))
	assert.Equal(t, fieldArgs["updated_at"].Type, graphql.DateTime)
}

func TestTranslator_MustTranslateArgs_ValueType_Panics(t *testing.T) {
	g := NewTranslator()

	assert.Panics(t, func() {
		g.MustTranslateArgs(2)
	})
}
