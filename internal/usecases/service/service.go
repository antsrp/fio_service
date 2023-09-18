package service

import (
	"github.com/antsrp/fio_service/internal/domain"
	"github.com/antsrp/fio_service/internal/interfaces/broker"
	"github.com/antsrp/fio_service/internal/interfaces/repository"
	iservice "github.com/antsrp/fio_service/internal/interfaces/service"
	"go.uber.org/zap"
)

var (
	_ iservice.Consumer = ServiceConsumer{}
	_ iservice.Producer = Service{}
)

type Service struct {
	Logger   *zap.SugaredLogger
	repo     repository.FIORepository
	producer broker.Producer[domain.PersonCommon]
}

func InitService(logger *zap.Logger, repo repository.FIORepository, producer broker.Producer[domain.PersonCommon]) *Service {
	return &Service{
		Logger:   logger.Sugar(),
		repo:     repo,
		producer: producer,
	}
}
