package main

import (
	"context"
	"log"

	"github.com/antsrp/fio_service/internal/domain"
	"github.com/antsrp/fio_service/internal/infrastructure/graphql"
	"github.com/antsrp/fio_service/internal/infrastructure/http/routes"
	"github.com/antsrp/fio_service/internal/infrastructure/http/server"
	"github.com/antsrp/fio_service/internal/infrastructure/http/web"
	"github.com/antsrp/fio_service/internal/usecases/service"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("Can't create zap logger: ", err)
	}

	configs, err := initConfigs()
	if err != nil {
		logger.Sugar().Fatal(err)
	}

	personStorage, err := initRepository(configs.rc, logger)
	if err != nil {
		logger.Sugar().Fatal(err)
	}
	defer handleCloser(logger, "person storage", personStorage)

	dataStorage, err := initDataStore(configs.dsc, logger)
	if err != nil {
		logger.Sugar().Fatal(err)
	}
	defer handleCloser(logger, "data storage", dataStorage)

	ctx, cancel := context.WithCancel(context.Background())
	entities := make(chan domain.PersonMessage)

	brokers, err := initBrokers(configs.bc, logger, ctx, entities)
	if err != nil {
		logger.Sugar().Fatal(err)
	}
	defer handleCloser(logger, "broker producer main", brokers.producerMain)
	defer handleCloser(logger, "broker producer aux", brokers.producerAux)
	defer handleCloser(logger, "broker consumer", brokers.consumerAux)

	webConnection := web.CreateNewConnection(logger)

	serviceMain, serviceConsumer := service.InitService(logger, personStorage, brokers.producerMain),
		service.NewServiceConsumer(logger, personStorage, brokers.producerAux, brokers.consumerAux, dataStorage,
			webConnection, entities, ctx)

	h := routes.NewHandler(serviceMain)
	api := graphql.NewPersonGQLApi(serviceMain)
	rts := h.Routes(api)

	end := make(chan bool)
	httpSrv := server.NewServer(configs.hc, logger, end, rts)
	logger.Info("server started")
	go serviceConsumer.Consume()
	go httpSrv.Start()

	<-end
	logger.Info("application to shutdown")
	cancel()
}
