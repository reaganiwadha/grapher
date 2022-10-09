package grapher

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/reaganiwadha/grapher/scalars"
	"reflect"
	"strings"
)

type OutputTranslationTable map[string]graphql.Output

var PrimitiveTranslationTable = OutputTranslationTable{
	"int":     graphql.Int,
	"int8":    graphql.Int,
	"int16":   graphql.Int,
	"int32":   graphql.Int,
	"int64":   graphql.Int,
	"uint":    graphql.Int,
	"uint8":   graphql.Int,
	"uint16":  graphql.Int,
	"uint32":  graphql.Int,
	"uint64":  graphql.Int,
	"float32": graphql.Float,
	"float64": graphql.Float,
	"string":  graphql.String,
	"bool":    graphql.Boolean,
}

// TranslatorConfig Stores the configuration of the graphql object builder
type TranslatorConfig struct {
	TypeTranslationTable     OutputTranslationTable
	DisablePointerToNullable bool
}

type translator struct {
	outputObjTable      OutputTranslationTable
	outputInputObjTable OutputTranslationTable
}

func getNamingByStructField(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")

	if jsonTag != "" {
		return strings.TrimSuffix(jsonTag, ",")
	}

	return field.Name
}

// New Returns a new translator
// It also stores already translated graphql.Object/graphql.InputObject to eliminate duplicates
func New() translator {
	return translator{
		outputObjTable:      OutputTranslationTable{},
		outputInputObjTable: OutputTranslationTable{},
	}
}

func (g translator) translateOutputRefType(t reflect.Type, inputObject bool) (ret graphql.Output, err error) {
	isPtr := false
	isArr := false

	if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		t = t.Elem()
		isArr = true
	}

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
		isPtr = true
	}

	if cached, ok := g.outputObjTable[t.Name()]; ok && !inputObject {
		ret = cached
	} else if cached, ok := g.outputInputObjTable[t.Name()]; ok && inputObject {
		ret = cached
	} else if prim, ok := PrimitiveTranslationTable[t.Name()]; ok {
		ret = prim
	} else if t.Kind() == reflect.Map {
		ret = scalars.ScalarJSON
	} else if t.Kind() == reflect.Struct {
		fields := graphql.Fields{}
		inputFields := graphql.InputObjectConfigFieldMap{}

		for i := 0; i < t.NumField(); i++ {
			fieldT := t.Field(i).Type
			fieldOut, _ := g.translateOutputRefType(fieldT, inputObject)

			if inputObject {
				inputFields[getNamingByStructField(t.Field(i))] = &graphql.InputObjectFieldConfig{
					Type: fieldOut,
				}
			} else {
				fields[getNamingByStructField(t.Field(i))] = &graphql.Field{
					Type: fieldOut,
				}
			}
		}

		if inputObject {
			ret = graphql.NewInputObject(graphql.InputObjectConfig{
				Name:   t.Name(),
				Fields: inputFields,
			})
			g.outputInputObjTable[t.Name()] = ret
		} else {
			ret = graphql.NewObject(graphql.ObjectConfig{
				Name:   t.Name(),
				Fields: fields,
			})
			g.outputObjTable[t.Name()] = ret
		}
	}

	if !isPtr {
		ret = graphql.NewNonNull(ret)
	}

	if isArr {
		ret = graphql.NewList(ret)
	}

	return
}

// Translate translates the type into a *graphql.Object
func (g translator) Translate(t interface{}) (ret graphql.Output, err error) {
	v := reflect.TypeOf(t)
	return g.translateOutputRefType(v, false)
}

// MustTranslate calls Translate, but will panic on error
func (g translator) MustTranslate(t interface{}) graphql.Output {
	ret, _ := g.Translate(t)

	// No errors as of now
	//if err != nil {
	//	panic(err)
	//}

	return ret
}

// TranslateInputObject Translate translates a struct into a *graphql.InputObject
func (g translator) TranslateInputObject(t interface{}) (*graphql.InputObject, error) {
	v := reflect.TypeOf(t)
	testV := v

	if testV.Kind() == reflect.Pointer {
		testV = testV.Elem()
	}

	if testV.Kind() != reflect.Struct {
		return nil, fmt.Errorf("grapher: TranslateInputObject only works with a struct type")
	}

	out, err := g.translateOutputRefType(v, true)

	return out.(*graphql.InputObject), err
}

//
//// MustTranslateInputObject calls TranslateInputObject, but will panic on error
//func (g translator) MustTranslateInputObject(t interface{}) *graphql.InputObject {
//	ret, err := g.TranslateInputObject(t)
//
//	if err != nil {
//		panic(err)
//	}
//
//	return ret
//}
