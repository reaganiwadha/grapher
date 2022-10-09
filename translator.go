package grapher

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/reaganiwadha/grapher/scalars"
	"reflect"
	"strings"
)

// OutputTranslationTable is a map that stores the graphql.Output according to it's type
type OutputTranslationTable map[string]graphql.Output

var primitiveTranslationTable = OutputTranslationTable{
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

// Translator purpose is to provide an easy way of translating various types into graphql types
type Translator interface {
	Translate(t interface{}) (ret graphql.Output, err error)
	TranslateInput(t interface{}) (ret *graphql.InputObject, err error)
	TranslateArgs(t interface{}) (ret graphql.FieldConfigArgument, err error)

	MustTranslate(t interface{}) graphql.Output
	MustTranslateInput(t interface{}) *graphql.InputObject
	MustTranslateArgs(t interface{}) (ret graphql.FieldConfigArgument)
}

// New Returns a new translator
// It also stores already translated graphql.Object/graphql.InputObject to eliminate duplicates
func New() Translator {
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
	} else if prim, ok := primitiveTranslationTable[t.Name()]; ok {
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

func assertStruct(v reflect.Type) (structType reflect.Type, err error) {
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("grapher: not a struct type")
	}

	return v, nil
}

// TranslateInput translates a struct into a *graphql.InputObject
func (g translator) TranslateInput(t interface{}) (ret *graphql.InputObject, err error) {
	v := reflect.TypeOf(t)
	if _, err = assertStruct(v); err != nil {
		return
	}

	out, err := g.translateOutputRefType(v, true)

	return out.(*graphql.InputObject), err
}

// MustTranslateInput calls TranslateInputObject, but will panic on error
func (g translator) MustTranslateInput(t interface{}) *graphql.InputObject {
	ret, err := g.TranslateInput(t)

	if err != nil {
		panic(err)
	}

	return ret
}

// TranslateArgs translates a struct into a graphql.FieldConfigArgument
func (g translator) TranslateArgs(t interface{}) (ret graphql.FieldConfigArgument, err error) {
	structType, err := assertStruct(reflect.TypeOf(t))

	if err != nil {
		return
	}

	ret = graphql.FieldConfigArgument{}

	for i := 0; i < structType.NumField(); i++ {
		field, _ := g.translateOutputRefType(structType.Field(i).Type, true)
		//if fieldErr != nil{
		//	return nil, fieldErr
		//}

		ret[getNamingByStructField(structType.Field(i))] = &graphql.ArgumentConfig{
			Type: field,
		}
	}

	return
}

// MustTranslateArgs calls TranslateArgs, but will panic on error
func (g translator) MustTranslateArgs(t interface{}) (ret graphql.FieldConfigArgument) {
	ret, err := g.TranslateArgs(t)

	if err != nil {
		panic(err)
	}

	return
}
