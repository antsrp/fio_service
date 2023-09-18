package repository

import "github.com/antsrp/fio_service/internal/domain"

type FIORepository interface {
	Add(domain.Person) *RepositoryError
	Get([]string, []string, int) ([]domain.Person, *RepositoryError)
	Update(domain.Person) (int, *RepositoryError)
	Delete(int) (int, *RepositoryError)
	DBConnection
}

type RepositoryError struct {
	IsInternal bool
	Cause      error
}
