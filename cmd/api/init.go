package main

import (
	"context"
	"fmt"

	"github.com/antsrp/fio_service/internal/config"
	"github.com/antsrp/fio_service/internal/domain"
	"github.com/antsrp/fio_service/internal/domain/benefication"
	"github.com/antsrp/fio_service/internal/infrastructure/broker"
	datastore "github.com/antsrp/fio_service/internal/infrastructure/data_store"
	"github.com/antsrp/fio_service/internal/infrastructure/http"
	repo "github.com/antsrp/fio_service/internal/infrastructure/repository"
	ibroker "github.com/antsrp/fio_service/internal/interfaces/broker"
	idstore "github.com/antsrp/fio_service/internal/interfaces/data_store"
	"github.com/antsrp/fio_service/internal/interfaces/repository"
	"github.com/antsrp/fio_service/internal/usecases/broker/kafka"
	"github.com/antsrp/fio_service/internal/usecases/data_store/redis"
	"github.com/antsrp/fio_service/internal/usecases/database/postgres"
	"go.uber.org/zap"
)

type Configs struct {
	rc  *repo.Settings
	dsc *datastore.Settings
	hc  *http.Settings
	bc  *broker.Settings
}

type Brokers struct {
	producerMain ibroker.Producer[domain.PersonCommon]
	producerAux  ibroker.Producer[domain.PersonMessage]
	consumerAux  ibroker.Consumer[domain.PersonMessage]
}

func initConfigs() (*Configs, error) {
	rc, err := config.Parse[repo.Settings]("DB")
	if err != nil {
		return nil, fmt.Errorf("can't launch config of repository: %w", err)
	}
	dsc, err := config.Parse[datastore.Settings]("STORE")
	if err != nil {
		return nil, fmt.Errorf("can't launch config of repository: %w", err)
	}
	hc, err := config.Parse[http.Settings]("SERVER")
	if err != nil {
		return nil, fmt.Errorf("can't launch config of repository: %w", err)
	}
	bc, err := config.Parse[broker.Settings]("BROKER")
	if err != nil {
		return nil, fmt.Errorf("can't launch config of repository: %w", err)
	}

	return &Configs{
		rc:  rc,
		dsc: dsc,
		hc:  hc,
		bc:  bc,
	}, nil
}

func initRepository(dc *repo.Settings, logger *zap.Logger) (repository.FIORepository, error) {
	var personStorage repository.FIORepository
	switch dc.Type {
	case "postgres":
		dbConn, err := postgres.NewConnection(dc, logger)
		if err != nil {
			return nil, fmt.Errorf("can't create postgres connection: %w", err)
		}
		personStorage, err = postgres.NewPersonStorage(dbConn, dc.CountOnPage)
		if err != nil {
			return nil, fmt.Errorf("can't create person storage: %w", err)
		}
	default:
		return nil, fmt.Errorf("there is no database type in .env file")
	}

	return personStorage, nil
}

func initDataStore(dsc *datastore.Settings, logger *zap.Logger) (idstore.Storager[benefication.DataCacheValue], error) {
	var dataStorage idstore.Storager[benefication.DataCacheValue]
	var err error
	switch dsc.Type {
	case "redis":
		dataStorage, err = redis.NewStorage[benefication.DataCacheValue](dsc, logger)
		if err != nil {
			return nil, fmt.Errorf("can't create connection to redis storage: %w", err)
		}
	default:
		return nil, fmt.Errorf("there is no data store type in .env file")
	}

	return dataStorage, nil
}

func initBrokers(bc *broker.Settings, logger *zap.Logger,
	ctx context.Context, entities chan domain.PersonMessage) (*Brokers, error) {

	var brokerProducerMain ibroker.Producer[domain.PersonCommon]
	var brokerProducerAux ibroker.Producer[domain.PersonMessage]
	var brokerConsumerAux ibroker.Consumer[domain.PersonMessage]

	var err error
	switch bc.Type {
	case "kafka":
		if brokerProducerMain, err = kafka.NewProducer[domain.PersonCommon](bc.Host, bc.Port, bc.Topic, logger); err != nil {
			return nil, fmt.Errorf("can't create kafka producer main: %w", err)
		}
		if brokerProducerAux, err = kafka.NewProducer[domain.PersonMessage](bc.Host, bc.Port, bc.TopicFailed, logger); err != nil {
			return nil, fmt.Errorf("can't create kafka producer aux: %w", err)
		}
		if brokerConsumerAux, err = kafka.NewConsumer[domain.PersonMessage](bc.Host, bc.Port, bc.Topic, bc.GroupID, logger, entities, ctx); err != nil {
			return nil, fmt.Errorf("can't create kafka consumer: %w", err)
		}
	default:
		return nil, fmt.Errorf("there is no broker type in .env file")
	}

	return &Brokers{
		producerMain: brokerProducerMain,
		producerAux:  brokerProducerAux,
		consumerAux:  brokerConsumerAux,
	}, nil
}
