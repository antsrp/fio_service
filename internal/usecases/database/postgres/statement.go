package postgres

import (
	"database/sql"

	"github.com/pkg/errors"
)

// StatementStorage struct contains struct of connection and statements
type StatementStorage struct {
	conn       *Connection
	statements []*sql.Stmt
}

// Create new stmt storage
func Create(c *Connection) StatementStorage {
	return StatementStorage{conn: c}
}

// Close implements io.Closer interface. It is used for close statements (graceful shutdown)
func (s *StatementStorage) Close() {
	for _, stmt := range s.statements {
		if err := stmt.Close(); err != nil {
			s.conn.logger.Sugar().Errorf("can't close statement: %v", err.Error())
		}
	}
}

type stmt struct {
	Query string
	Dst   **sql.Stmt
}

func (s *StatementStorage) prepareStatement(query string) (*sql.Stmt, error) {
	stmt, err := s.conn.db.Prepare(query)
	if err != nil {
		return nil, errors.Wrapf(err, "can't prepare query %q", query)
	}

	return stmt, nil
}

func (s *StatementStorage) initStatements(statements []stmt) error {
	for i := range statements {
		statement, err := s.prepareStatement(statements[i].Query)
		if err != nil {
			return err
		}

		*statements[i].Dst = statement
		s.statements = append(s.statements, statement)
	}

	return nil
}
