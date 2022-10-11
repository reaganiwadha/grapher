# grapher 
[![Go Report Card](https://goreportcard.com/badge/github.com/reaganiwadha/grapher)](https://goreportcard.com/report/github.com/reaganiwadha/grapher) [![codecov](https://codecov.io/gh/reaganiwadha/grapher/branch/trunk/graph/badge.svg)](https://codecov.io/gh/reaganiwadha/grapher)
Neat extra tooling for [graphql-go](https://github.com/graphql-go/graphql).

## v1 Development
This package is still a work-in-progress. The feature candidates for the v1 release is :
- [x] Struct to `graphql.Object`/`graphql.InputObject` translator utilizing struct tags
- [x] Struct to `map[string]*graphql.ArgumentConfig{}`
- [x] Configurable struct tags translator
- [x] Schema field builder utilizing go 1.18 Generics feature ensuring type safety
- [ ] Resolver middlewares
- [x] Mux-like schema building with `MutationQueryCollector`

If you have any feature ideas, do not hesitate to tell and open a new issue in the Issues tab!

## Contributing
Contributions in any way, shape, or form is super welcomed! If you have any issues, ideas or even a question just open a new issue, or make a pull request. 

## License
MIT
