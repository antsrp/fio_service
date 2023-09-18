package graphql

import (
	"fmt"

	"github.com/antsrp/fio_service/internal/domain"
	igql "github.com/antsrp/fio_service/internal/interfaces/graphql"
	"github.com/graphql-go/graphql"
)

var personType = graphql.NewObject(graphql.ObjectConfig{
	Name:   "Person",
	Fields: graphql.BindFields(domain.Person{}),
})

func initQuery(get func(params graphql.ResolveParams) (interface{}, error)) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "PersonQuery",
		Fields: graphql.Fields{
			"getPersons": &graphql.Field{
				Type:        graphql.NewList(personType),
				Description: "Get persons with filters",
				Args: graphql.FieldConfigArgument{
					"where": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"order": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"page": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: get,
			},
		},
	})
}

func initMutations(api igql.PersonAPI) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "PersonMutation",
		Fields: graphql.Fields{
			"addPerson": &graphql.Field{
				Type:        graphql.String,
				Description: "add new person",
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"surname": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"patronym": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: api.Add,
			},
			"updatePerson": &graphql.Field{
				Type:        graphql.String,
				Description: "update existing person by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"surname": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"patronym": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"age": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"gender": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"nationality": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: api.Update,
			},
			"deletePerson": &graphql.Field{
				Type:        graphql.String,
				Description: "delete existing person",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: api.Delete,
			},
		},
	})
}

func InitSchema(api igql.PersonAPI) *graphql.Schema {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    initQuery(api.Get),
		Mutation: initMutations(api),
	})
	if err != nil {
		fmt.Printf("error with schema: %v\n", err.Error())
	}

	return &schema
}
