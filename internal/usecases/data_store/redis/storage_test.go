package redis

import (
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	"github.com/antsrp/fio_service/internal/config"
	"github.com/antsrp/fio_service/internal/domain/benefication"
	datastore "github.com/antsrp/fio_service/internal/infrastructure/data_store"
	idstore "github.com/antsrp/fio_service/internal/interfaces/data_store"
	"go.uber.org/zap"
)

var dataStorage idstore.Storager[benefication.DataCacheValue]

func TestMain(m *testing.M) {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("Can't create zap logger: ", err)
	}

	dsc, err := config.Parse[datastore.Settings]("STORE")
	if err != nil {
		logger.Sugar().Fatal("can't launch config of repository: %v", err.Error())
	}
	dsc.DBName = 3
	dsc.ExpirationTime = 1

	dataStorage, err = NewStorage[benefication.DataCacheValue](dsc, logger)
	if err != nil {
		logger.Sugar().Fatal("can't create connection to redis storage: %v", err.Error())
	}
	defer dataStorage.Close()

	m.Run()
}

var (
	name  = "Alexey"
	value = benefication.DataCacheValue{
		Age:    46,
		Gender: "male",
		Country: []benefication.Country{
			{ID: "RU", Probability: 0.445},
			{ID: "UA", Probability: 0.213},
			// etc
		},
	}
)

func TestGetWithoutSet(t *testing.T) {
	val, err := dataStorage.Get(name)

	assert.EqualErrorf(t, err, "redis: nil", "expected error: %v, got %v", redis.Nil, err.Error())
	assert.Equal(t, benefication.DataCacheValue{}, val, "there is no expected value, got %v", val)
}

func TestSet(t *testing.T) {
	err := dataStorage.Set(name, value)

	assert.NoError(t, err, "there should be no error when executing this request")
}

func TestGet(t *testing.T) {
	val, err := dataStorage.Get(name)

	assert.NoError(t, err, "there should be no error when executing this request")
	assert.Equal(t, value, val, "name %s, expected value %v, got %v", name, value, val)
}

func TestDelete(t *testing.T) {
	err := dataStorage.Delete(name)

	assert.NoError(t, err, "there should be no error when executing this request")
}

func TestGetDeleted(t *testing.T) {
	val, err := dataStorage.Get(name)

	assert.EqualErrorf(t, err, "redis: nil", "expected error: %v, got %v", redis.Nil, err.Error())
	assert.Equal(t, benefication.DataCacheValue{}, val, "there is no expected value, got %v", val)
}

func TestGetNoKey(t *testing.T) {
	val, err := dataStorage.Get("Sergey")

	assert.EqualErrorf(t, err, "redis: nil", "expected error: %v, got %v", redis.Nil, err.Error())
	assert.Equal(t, benefication.DataCacheValue{}, val, "there is no expected value, got %v", val)
}
