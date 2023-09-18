package service

import (
	"context"

	"github.com/antsrp/fio_service/internal/domain"
	"github.com/antsrp/fio_service/internal/domain/benefication"
	"github.com/antsrp/fio_service/internal/interfaces/broker"
	idstore "github.com/antsrp/fio_service/internal/interfaces/data_store"
	irepo "github.com/antsrp/fio_service/internal/interfaces/repository"
	iservice "github.com/antsrp/fio_service/internal/interfaces/service"
	"github.com/antsrp/fio_service/internal/interfaces/web"
	"go.uber.org/zap"
)

var _ iservice.Consumer = ServiceConsumer{}

type ServiceConsumer struct {
	Logger   *zap.SugaredLogger
	repo     irepo.FIORepository
	producer broker.Producer[domain.PersonMessage]
	consumer broker.Consumer[domain.PersonCommon]
	dstorage idstore.Storager[benefication.DataCacheValue]
	conn     web.Connector
	channel  chan domain.PersonMessage
	ctx      context.Context
}

func NewServiceConsumer(logger *zap.Logger, repo irepo.FIORepository, producer broker.Producer[domain.PersonMessage],
	consumer broker.Consumer[domain.PersonCommon], ds idstore.Storager[benefication.DataCacheValue], conn web.Connector,
	channel chan domain.PersonMessage, ctx context.Context) *ServiceConsumer {
	return &ServiceConsumer{
		Logger:   logger.Sugar(),
		repo:     repo,
		producer: producer,
		consumer: consumer,
		dstorage: ds,
		conn:     conn,
		channel:  channel,
		ctx:      ctx,
	}
}

func (s ServiceConsumer) Consume() {
	go s.consumer.Consume()

	var msg domain.PersonMessage
	for s.ctx.Err() == nil {
		msg = <-s.channel
		s.Logger.Infof("message consumed, name: %s, surname: %s, patronym: %s", msg.Name, msg.Surname, msg.Patronym)
		if msg.Name == "" {
			msg.EntityError = "person name wasn't set"
		} else if msg.Surname == "" {
			msg.EntityError = "person surname wasn't set"
		}
		if msg.EntityError != "" {
			if err := s.producer.Produce("fail", msg); err != nil {
				s.Logger.Errorf("can't produce message to topic: %v", err.Error())
			}
			continue
		}
		person := domain.Person{
			PersonCommon: domain.PersonCommon{
				Name:     msg.Name,
				Surname:  msg.Surname,
				Patronym: msg.Patronym,
			},
		}
		val, err := s.dstorage.Get(msg.Name)
		if err != nil {
			s.Logger.Infof("can't find value in cache for name %s, have to make requests", msg.Name)
			if val, err = beneficate(s.conn, msg.Name); err != nil {
				s.Logger.Errorf("can't beneficate person: %v", err.Error())
				break
			}
			if err := s.dstorage.Set(msg.Name, val); err != nil {
				s.Logger.Errorf("can't add benefication data to store: %v", err.Error())
			}
		}
		s.Logger.Infof("beneficated value: age: %d, gender: %s, countries: %v", val.Age, val.Gender, val.Country)
		person.Age = val.Age
		person.Gender = val.Gender
		if len(val.Country) > 0 {
			person.Nationality = val.Country[0].ID
		}

		if err := s.repo.Add(person); err != nil {
			s.Logger.Errorf("can't add person to db: %v", err.Cause.Error())
		}
	}
}
