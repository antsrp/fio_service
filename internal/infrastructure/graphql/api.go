package graphql

import (
	"fmt"
	"strings"

	"github.com/antsrp/fio_service/internal/domain"
	"github.com/antsrp/fio_service/internal/infrastructure/messages"
	igql "github.com/antsrp/fio_service/internal/interfaces/graphql"
	"github.com/antsrp/fio_service/internal/usecases/service"
	"github.com/graphql-go/graphql"
)

type PersonGQLApi struct {
	srv *service.Service
}

func NewPersonGQLApi(srv *service.Service) *PersonGQLApi {
	return &PersonGQLApi{
		srv: srv,
	}
}

func (p PersonGQLApi) Add(params graphql.ResolveParams) (interface{}, error) {
	name, _ := params.Args["name"].(string)
	surname, _ := params.Args["surname"].(string)
	patronym, _ := params.Args["patronym"].(string)

	person := domain.PersonCommon{
		Name:     name,
		Surname:  surname,
		Patronym: patronym,
	}

	if err := p.srv.SendToBroker(person); err != nil {
		p.srv.Logger.Errorf("can't send person to broker via graphql query: %v", err.Error())
		return nil, err
	}
	return messages.SendToAddMsg, nil
}

func (p PersonGQLApi) Get(params graphql.ResolveParams) (interface{}, error) {
	where, _ := params.Args["where"].(string)
	order, _ := params.Args["order"].(string)
	page, _ := params.Args["page"].(int)

	var wheres, orders []string
	if where != "" {
		wheres = strings.Split(where, ",")
	}
	if order != "" {
		orders = strings.Split(order, ",")
	}

	persons, customErr := p.srv.GetPersons(wheres, orders, page)
	if customErr != nil {
		p.srv.Logger.Errorf("can't get persons via graphql query: %v", customErr.Cause.Error())
		var err error
		if customErr.IsInternal {
			err = fmt.Errorf(messages.InternalError)
		} else {
			err = fmt.Errorf(messages.InvalidInput)
		}
		return nil, err
	}
	return persons, nil
}

func (p PersonGQLApi) Update(params graphql.ResolveParams) (interface{}, error) {
	id, ok := params.Args["id"].(int)
	if !ok {
		return nil, fmt.Errorf(fmt.Sprintf(messages.IncorrectAssignment, "id"))
	}
	name, _ := params.Args["name"].(string)
	surname, _ := params.Args["surname"].(string)
	patronym, _ := params.Args["patronym"].(string)
	age, _ := params.Args["age"].(int)
	gender, _ := params.Args["gender"].(string)
	nationality, _ := params.Args["nationality"].(string)
	person := domain.Person{
		Id: id,
		PersonCommon: domain.PersonCommon{
			Name:     name,
			Surname:  surname,
			Patronym: patronym,
		},
		Age:         uint(age),
		Gender:      gender,
		Nationality: nationality,
	}

	rowsAffected, customErr := p.srv.UpdatePerson(person)
	if customErr != nil {
		p.srv.Logger.Errorf("can't update person via graphql query: %v", customErr.Cause.Error())
		var err error
		if customErr.IsInternal {
			err = fmt.Errorf(messages.InternalError)
		} else {
			err = fmt.Errorf(messages.InvalidInput)
		}
		return nil, err
	}
	if rowsAffected == 0 {
		return messages.NoRowsAffected, nil
	}
	return fmt.Sprintf(messages.SuccessfulUpdateMsg, id), nil
}

func (p PersonGQLApi) Delete(params graphql.ResolveParams) (interface{}, error) {
	id, ok := params.Args["id"].(int)
	if !ok {
		return nil, fmt.Errorf(fmt.Sprintf(messages.IncorrectAssignment, "id"))
	}
	rowsAffected, customErr := p.srv.DeletePerson(id)
	if customErr != nil {
		p.srv.Logger.Errorf("can't delete person via graphql query: %v", customErr.Cause.Error())
		var err error
		if customErr.IsInternal {
			err = fmt.Errorf(messages.InternalError)
		} else {
			err = fmt.Errorf(messages.InvalidInput)
		}
		return nil, err
	}
	if rowsAffected == 0 {
		return messages.NoRowsAffected, nil
	}
	return fmt.Sprintf(messages.SuccessfulDeleteMsg, id), nil
}

var _ igql.PersonAPI = PersonGQLApi{}
