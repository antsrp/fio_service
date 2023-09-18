package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/antsrp/fio_service/internal/domain/benefication"
	datastore "github.com/antsrp/fio_service/internal/infrastructure/data_store"
	idstore "github.com/antsrp/fio_service/internal/interfaces/data_store"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Storage[T any] struct {
	client     *redis.Client
	logger     *zap.Logger
	expiration time.Duration
	ctx        context.Context
}

func NewStorage[T any](settings *datastore.Settings, l *zap.Logger) (*Storage[T], error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", settings.Host, settings.Port),
		Password: settings.Password,
		DB:       settings.DBName,
	})

	ctx := context.Background()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}
	l.Info("redis connection opened")

	return &Storage[T]{
		client:     client,
		logger:     l,
		expiration: time.Duration(settings.ExpirationTime) * time.Minute,
		ctx:        ctx,
	}, nil
}

func (s Storage[T]) Set(key string, value T) error {
	return s.client.Set(s.ctx, key, value, s.expiration).Err()
}

func (s Storage[T]) Get(key string) (T, error) {
	status := s.client.Get(s.ctx, key)
	if err := status.Err(); err != nil {
		return *new(T), err
	}
	var returnVal T
	if err := status.Scan(&returnVal); err != nil {
		return *new(T), err
	}
	return returnVal, nil
}

func (s Storage[T]) Delete(key string) error {
	return s.client.Del(s.ctx, key).Err()
}

func (s Storage[T]) Close() error {
	s.logger.Info("redis connection closing")
	err := s.client.Close()
	if err != nil {
		s.logger.Sugar().Errorf("can't close connection to redis data storage: %v", err.Error())
	} else {
		s.logger.Info("redis connection closed")
	}
	return err
}

var _ idstore.Storager[benefication.DataCacheValue] = Storage[benefication.DataCacheValue]{}
