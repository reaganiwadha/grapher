package grapher

import (
	"encoding/json"
	"github.com/graphql-go/graphql"
	"reflect"
)

type ResolverFunc[argsT any, outT any] func(p graphql.ResolveParams, args argsT) (outT, error)

type GuardFunc func(p graphql.ResolveParams) bool

type NoArgs struct{}

type FieldBuilderConfig struct {
	Translator Translator
}

// FieldBuilder interface, put NoArgs to use no args
type FieldBuilder[argT any, outT any] interface {
	WithDescription(desc string) FieldBuilder[argT, outT]
	WithResolver(resolverFunc ResolverFunc[argT, outT]) FieldBuilder[argT, outT]

	Build() (field *graphql.Field, err error)

	// TODO
	//AddGuard(guardFunc GuardFunc) FieldBuilder[argT, outT]
}

type fieldBuilder[argT any, outT any] struct {
	desc     string
	resolver ResolverFunc[argT, outT]

	translator Translator
}

func (f *fieldBuilder[argT, outT]) WithDescription(desc string) FieldBuilder[argT, outT] {
	f.desc = desc
	return f
}

func (f *fieldBuilder[argT, outT]) WithResolver(resolverFunc ResolverFunc[argT, outT]) FieldBuilder[argT, outT] {
	f.resolver = resolverFunc
	return f
}

func (f *fieldBuilder[argT, outT]) Build() (field *graphql.Field, err error) {
	field = &graphql.Field{}
	field.Description = f.desc

	useArgs := !isGrapherNoArgs[argT]()

	if useArgs {
		fieldArgT, argErr := f.translator.TranslateArgs(new(argT))
		if argErr != nil {
			err = argErr
			return
		}

		field.Args = fieldArgT
	}

	fieldOutT, _ := f.translator.Translate(new(outT))
	//if err != nil {
	//	return
	//}

	field.Resolve = createResolver[argT, outT](useArgs, f.resolver)
	field.Type = fieldOutT

	return
}

func NewFieldBuilder[argT any, outT any](cfgArgs ...FieldBuilderConfig) FieldBuilder[argT, outT] {
	if len(cfgArgs) != 0 {
		if cfgArgs[0].Translator != nil {
			return &fieldBuilder[argT, outT]{
				translator: cfgArgs[0].Translator,
			}
		}

	}

	return &fieldBuilder[argT, outT]{
		translator: NewTranslator(),
	}
}

func isGrapherNoArgs[T any]() bool {
	return reflect.TypeOf(new(T)).String() == "*grapher.NoArgs"
}

func createResolver[aT any, oT any](useArgs bool, mainResolver func(p graphql.ResolveParams, args aT) (oT, error)) func(p graphql.ResolveParams) (res interface{}, err error) {
	return func(p graphql.ResolveParams) (res interface{}, err error) {
		args := new(aT)

		if !useArgs {
			return mainResolver(p, *args)
		}

		r, err := json.Marshal(p.Args)

		if err != nil {
			return
		}

		json.Unmarshal(r, &args)

		return mainResolver(p, *args)
	}
}
