package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/antsrp/fio_service/internal/domain"
	irepo "github.com/antsrp/fio_service/internal/interfaces/repository"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type PersonStorage struct {
	limit int
	StatementStorage

	createPersonStmt *sql.Stmt
	deletePersonStmt *sql.Stmt
}

var _ irepo.FIORepository = &PersonStorage{}

const (
	createPersonQ = `INSERT INTO persons (name, surname, patronym, age, gender, nationality) VALUES 
	($1, $2, $3, $4, $5, $6) RETURNING id`
	deletePersonQ = "DELETE FROM persons WHERE id = $1"
	updateMainQ   = "UPDATE persons SET"
	selectMainQ   = "SELECT id, name, surname, patronym, age, gender, nationality FROM Persons"
	where         = "WHERE"
	orderby       = "ORDER BY"
	limit         = "LIMIT"
	offset        = "OFFSET"

	classInternalError = "internal_error"
)

func NewPersonStorage(conn *Connection, limit int) (*PersonStorage, error) {
	s := &PersonStorage{
		StatementStorage: Create(conn),
		limit:            limit,
	}

	stmts := []stmt{
		{Query: createPersonQ, Dst: &s.createPersonStmt},
		{Query: deletePersonQ, Dst: &s.deletePersonStmt},
	}

	if err := s.initStatements(stmts); err != nil {
		return nil, errors.Wrap(err, "can't init statements")
	}

	return s, nil
}

func (ps PersonStorage) Add(person domain.Person) *irepo.RepositoryError {
	var id int
	if err := ps.createPersonStmt.
		QueryRow(&person.Name, &person.Surname, &person.Patronym, &person.Age, &person.Gender, &person.Nationality).
		Scan(&id); err != nil {
		var repoError irepo.RepositoryError
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code.Class().Name() == classInternalError {
				repoError.IsInternal = true
			}
		} else {
			if err := ps.Check(); err != nil {
				repoError.IsInternal = true
			}
		}
		repoError.Cause = errors.Wrap(err, "can't add person into database")
		return &repoError
	}
	return nil
}

func (ps PersonStorage) Delete(id int) (int, *irepo.RepositoryError) {
	result, err := ps.deletePersonStmt.Exec(&id)
	if err != nil {
		return 0, &irepo.RepositoryError{
			IsInternal: true,
			Cause:      errors.Wrapf(err, "can't delete person with id %d from database", id),
		}
	}
	rows, err := result.RowsAffected()
	if err != nil {
		ps.conn.logger.Sugar().Infof("can't get count of rows affected: %v", err.Error())
	}
	return int(rows), nil
}

func (ps PersonStorage) Get(wheres, order []string, page int) ([]domain.Person, *irepo.RepositoryError) {
	var builder strings.Builder
	builder.WriteString(selectMainQ)

	delim := " "

	if len(wheres) > 0 {
		builder.WriteString(delim)
		builder.WriteString(where)
		builder.WriteString(delim)
		builder.WriteString(strings.Join(wheres, " "))
	}
	if len(order) > 0 {
		builder.WriteString(delim)
		builder.WriteString(orderby)
		builder.WriteString(delim)
		builder.WriteString(strings.Join(order, ","))
	}

	if page > 0 {
		builder.WriteString(fmt.Sprintf("%s%s %d %s %d", delim, limit, ps.limit, offset, (page-1)*ps.limit))
	}

	cmd := builder.String()
	rows, err := ps.conn.db.Query(cmd)
	if err != nil {
		var repoError irepo.RepositoryError
		pqErr, ok := err.(*pq.Error)
		if ok {
			if pqErr.Code.Class().Name() == classInternalError {
				repoError.IsInternal = true
			}
		} else {
			if err := ps.Check(); err != nil {
				repoError.IsInternal = true
			}
		}
		repoError.Cause = errors.Wrapf(err, "can't get persons from database")
		return nil, &repoError
	}
	defer rows.Close()
	var persons []domain.Person
	for rows.Next() {
		var person domain.Person
		if err := rows.Scan(&person.Id, &person.Name, &person.Surname, &person.Patronym, &person.Age, &person.Gender, &person.Nationality); err != nil {
			return nil, &irepo.RepositoryError{
				IsInternal: true,
				Cause:      errors.Wrapf(err, "can't scan person into structure"),
			}
		}
		persons = append(persons, person)
	}

	return persons, nil
}

func (ps PersonStorage) Update(person domain.Person) (int, *irepo.RepositoryError) {
	var builder strings.Builder
	builder.WriteString(updateMainQ)

	delim := " "

	if person.Name != "" {
		builder.WriteString(fmt.Sprintf("%sname = '%s'", delim, person.Name))
		delim = ", "
	}
	if person.Surname != "" {
		builder.WriteString(fmt.Sprintf("%ssurname = '%s'", delim, person.Surname))
		delim = ", "
	}
	if person.Patronym != "" {
		builder.WriteString(fmt.Sprintf("%spatronym = '%s'", delim, person.Patronym))
		delim = ", "
	}
	if person.Age != 0 {
		builder.WriteString(fmt.Sprintf("%sage = %d", delim, person.Age))
		delim = ", "
	}
	if person.Gender != "" {
		builder.WriteString(fmt.Sprintf("%sgender = '%s'", delim, person.Gender))
		delim = ", "
	}
	if person.Nationality != "" {
		builder.WriteString(fmt.Sprintf("%snationality = '%s'", delim, person.Nationality))
	}

	builder.WriteString(fmt.Sprintf(" %s id = %d", where, person.Id))

	cmd := builder.String()
	result, err := ps.conn.db.Exec(cmd)
	if err != nil {
		return 0, &irepo.RepositoryError{
			IsInternal: true,
			Cause:      errors.Wrapf(err, "can't update person with id %d", person.Id),
		}
	}
	rows, err := result.RowsAffected()
	if err != nil {
		ps.conn.logger.Sugar().Infof("can't get count of rows affected: %v", err.Error())
	}
	return int(rows), nil
}

func (ps PersonStorage) Check() error {
	return ps.conn.Check()
}

func (ps PersonStorage) Close() error {
	return ps.conn.Close()
}
