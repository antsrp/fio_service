package postgres

import (
	"database/sql"
	"fmt"

	repo "github.com/antsrp/fio_service/internal/infrastructure/repository"
	irepo "github.com/antsrp/fio_service/internal/interfaces/repository"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	_ "github.com/lib/pq" // postgres driver
)

type Connection struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewConnection(settings *repo.Settings, logger *zap.Logger) (*Connection, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		settings.User, settings.Password, settings.Host, settings.Port, settings.DBName)

	db, err := sql.Open(settings.Type, dsn)

	if err != nil {
		return nil, errors.Wrap(err, "can't connect to db")
	}
	c := Connection{
		db:     db,
		logger: logger,
	}

	if err := c.Check(); err != nil {
		return nil, err
	}
	c.logger.Info("postgres connection opened")

	return &c, nil
}

func (c Connection) Check() error {
	return c.db.Ping()
}

func (c Connection) Close() error {
	c.logger.Info("postgres connection closing")
	err := c.db.Close()
	if err != nil {
		c.logger.Sugar().Errorf("can't close db connection: %v\n", err.Error())
	} else {
		c.logger.Info("postgres connection closed")
	}
	return err
}

var _ irepo.DBConnection = Connection{}
