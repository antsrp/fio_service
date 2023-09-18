package service

import (
	"github.com/antsrp/fio_service/internal/domain"
	repo "github.com/antsrp/fio_service/internal/interfaces/repository"
)

type Producer interface {
	SendToBroker(domain.PersonCommon) error
	GetPersons(wheres []string, order []string, page int) ([]domain.Person, *repo.RepositoryError)
	UpdatePerson(domain.Person) (int, *repo.RepositoryError)
	DeletePerson(int) (int, *repo.RepositoryError)
}

type Consumer interface {
	Consume()
}

type Multi interface {
	Producer
	Consumer
}
