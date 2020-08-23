package graphquery

import (
	"github.com/graphql-go/graphql"
	"fmt"
)
//Schema contains the schemas
var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: rootQuery,
	Mutation: rootMutation,
})

//ExecuteQuery handles graphql queries
func ExecuteQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema: schema,
		RequestString : query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}