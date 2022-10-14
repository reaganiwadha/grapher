# grapher [![Go Reference](https://pkg.go.dev/badge/github.com/graphql-go/graphql.svg)](https://pkg.go.dev/github.com/graphql-go/graphql) [![Go Report Card](https://goreportcard.com/badge/github.com/reaganiwadha/grapher)](https://goreportcard.com/report/github.com/reaganiwadha/grapher) [![codecov](https://codecov.io/gh/reaganiwadha/grapher/branch/trunk/graph/badge.svg)](https://codecov.io/gh/reaganiwadha/grapher)

Neat extra tooling for [graphql-go](https://github.com/graphql-go/graphql).

## Examples
### Without grapher
```go
type PostsQuery struct {
	Limit int         `json:"limit"`
	Query null.String `json:"query"`
	Tags  []string    `json:"tags"`
}

type Post struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func buildField() graphql.Fields {
	return graphql.Fields{
		"GetPosts": &graphql.Field{
			Description : "Returns posts",
			Type: graphql.NewList(graphql.NewObject(
				graphql.ObjectConfig{
					Name:        "Post",
					Fields:      graphql.Fields{
						"id" : &graphql.Field{
							Type: graphql.Int,
						},
						"body" : &graphql.Field{
							Type: graphql.String,
						},
					},
				},
			)),
			Args: map[string]*graphql.ArgumentConfig{
				"limit": {
					Type: graphql.NewNonNull(graphql.Int),
				},
				"query": {
					Type: graphql.String,
				},
				"tags": {
					Type: graphql.NewList(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (ret interface{}, err error) {
				var q PostsQuery
				if err = mapstructure.Decode(p.Args, &q); err != nil {
					return
				}

				/// resolver logic
				return
			},
		},
	}
}
```

### With grapher
```go
type PostsQuery struct {
	Limit int         `json:"limit"`
	Query null.String `json:"query"`
	Tags  []string    `json:"tags"`
}

type Post struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func buildFieldGrapher() graphql.Fields{
	return graphql.Fields{
		"GetPosts" : grapher.NewFieldBuilder[PostsQuery, []Post]().
			WithDescription("Returns posts").
			WithResolver(func(p graphql.ResolveParams, query PostsQuery) (ret []Post, err error) {
				// resolver logic
				return 
			}).
			MustBuild(),
	}
}
```

## v1 Development
This package is still a work-in-progress. The feature candidates for the v1 release is :
- [x] Struct to `graphql.Object`/`graphql.InputObject` translator utilizing struct tags
- [x] Struct to `map[string]*graphql.ArgumentConfig{}`
- [x] Configurable struct tags translator
- [x] Schema field builder utilizing go 1.18 Generics feature ensuring type safety
- [x] Resolver middlewares
- [x] Resolver arg validator
- [x] Mux-like schema building with `MutationQueryCollector`

Most features are done, but needs a little bit of polishing.

If you have any feature ideas, do not hesitate to tell and open a new issue in the Issues tab!

## Contributing
Contributions in any way, shape, or form is super welcomed! If you have any issues, ideas or even a question just open a new issue, or make a pull request. 

## License
MIT
