package grapher

import (
	"github.com/graphql-go/graphql"
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
	outputTranslationMap OutputTranslationTable
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
		outputTranslationMap: OutputTranslationTable{},
	}
}

func (g translator) translateOutputRefType(t reflect.Type) (ret graphql.Output, err error) {
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

	if cached, ok := g.outputTranslationMap[t.Name()]; ok {
		ret = cached
	} else if prim, ok := PrimitiveTranslationTable[t.Name()]; ok {
		ret = prim
	} else if t.Kind() == reflect.Map {
		ret = ScalarJSON
	} else if t.Kind() == reflect.Struct {
		fields := graphql.Fields{}

		for i := 0; i < t.NumField(); i++ {
			fieldT := t.Field(i).Type
			fieldOut, _ := g.translateOutputRefType(fieldT)

			fields[getNamingByStructField(t.Field(i))] = &graphql.Field{
				Type: fieldOut,
			}
		}

		ret = graphql.NewObject(graphql.ObjectConfig{
			Name:   t.Name(),
			Fields: fields,
		})

		g.outputTranslationMap[t.Name()] = ret
	}

	if !isPtr {
		ret = graphql.NewNonNull(ret)
	}

	if isArr {
		ret = graphql.NewList(ret)
	}

	return
}

// TranslateOutput Translates the type into a *graphql.Object
func (g translator) TranslateOutput(t interface{}) (ret graphql.Output, err error) {
	v := reflect.TypeOf(t)
	return g.translateOutputRefType(v)
}

// MustTranslateOutput calls TranslateOutput, but will panic on error
func (g translator) MustTranslateOutput(t interface{}) graphql.Output {
	ret, _ := g.TranslateOutput(t)

	// No errors as of now
	//if err != nil {
	//	panic(err)
	//}

	return ret
}

//func (g translator) TranslateInputObject(t interface{}) (*graphql.InputObject, error) {
//	return nil, nil
//}
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
