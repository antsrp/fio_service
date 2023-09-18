package graphql

import "github.com/graphql-go/graphql"

type PersonAPI interface {
	Add(params graphql.ResolveParams) (interface{}, error)
	Get(params graphql.ResolveParams) (interface{}, error)
	Update(params graphql.ResolveParams) (interface{}, error)
	Delete(params graphql.ResolveParams) (interface{}, error)
}
