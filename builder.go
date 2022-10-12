package grapher

import (
	"encoding/json"
	"github.com/graphql-go/graphql"
	"reflect"
)

type ResolverFn[argT any, outT any] func(p graphql.ResolveParams, arg argT) (outT, error)

type ResolverMiddlewareFn func(nextFn graphql.FieldResolveFn) graphql.FieldResolveFn

type ArgValidatorFn func(arg any) (err error)

type NoArgs struct{}

type FieldBuilderConfig struct {
	Translator Translator
}

// FieldBuilder interface, put NoArgs to use no args
type FieldBuilder[argT any, outT any] interface {
	WithDescription(desc string) FieldBuilder[argT, outT]
	WithResolver(resolverFunc ResolverFn[argT, outT]) FieldBuilder[argT, outT]
	WithCustomArgValidator(fn ArgValidatorFn) FieldBuilder[argT, outT]

	Build() (field *graphql.Field, err error)

	AddMiddleware(middlewareFn ResolverMiddlewareFn) FieldBuilder[argT, outT]
}

type fieldBuilder[argT any, outT any] struct {
	desc         string
	resolver     ResolverFn[argT, outT]
	argValidator ArgValidatorFn

	translator  Translator
	middlewares []ResolverMiddlewareFn
}

func (f *fieldBuilder[argT, outT]) WithDescription(desc string) FieldBuilder[argT, outT] {
	f.desc = desc
	return f
}

func (f *fieldBuilder[argT, outT]) WithResolver(resolverFunc ResolverFn[argT, outT]) FieldBuilder[argT, outT] {
	f.resolver = resolverFunc
	return f
}

func (f *fieldBuilder[argT, outT]) WithCustomArgValidator(fn ArgValidatorFn) FieldBuilder[argT, outT] {
	f.argValidator = fn
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

	field.Resolve = f.createResolver(useArgs)
	field.Type = fieldOutT

	return
}

func (f *fieldBuilder[argT, outT]) AddMiddleware(middlewareFn ResolverMiddlewareFn) FieldBuilder[argT, outT] {
	f.middlewares = append(f.middlewares, middlewareFn)

	return f
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

func (f *fieldBuilder[argT, outT]) createResolver(useArgs bool) func(p graphql.ResolveParams) (res interface{}, err error) {
	lastFn := func(p graphql.ResolveParams) (res interface{}, err error) {
		args := new(argT)

		if !useArgs {
			return f.resolver(p, *args)
		}

		r, _ := json.Marshal(p.Args)
		json.Unmarshal(r, &args)

		if f.argValidator != nil {
			if err = f.argValidator(args); err != nil {
				return
			}
		}

		return f.resolver(p, *args)
	}

	if len(f.middlewares) == 0 {
		return lastFn
	}

	endFn := f.middlewares[len(f.middlewares)-1](lastFn)

	for i := 1; i < len(f.middlewares); i++ {
		endFn = f.middlewares[i](endFn)
	}

	return endFn
}
