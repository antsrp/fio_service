package service

import (
	"strings"

	"github.com/antsrp/fio_service/internal/domain"
	irepo "github.com/antsrp/fio_service/internal/interfaces/repository"
)

func (s Service) SendToBroker(person domain.PersonCommon) error {
	return s.producer.Produce("", person)
}

func (s Service) GetPersons(wheres []string, order []string, page int) ([]domain.Person, *irepo.RepositoryError) {

	n := len(wheres)

	for i := 0; i < n; i++ {
		wheres[i] = strings.ReplaceAll(wheres[i], "|", " OR ")
		wheres[i] = strings.ReplaceAll(wheres[i], "!", " NOT ")
		wheres[i] = strings.ReplaceAll(wheres[i], "&", " AND ")
	}

	return s.repo.Get(wheres, order, page)
}

func (s Service) UpdatePerson(person domain.Person) (int, *irepo.RepositoryError) {
	return s.repo.Update(person)
}

func (s Service) DeletePerson(id int) (int, *irepo.RepositoryError) {
	return s.repo.Delete(id)
}
