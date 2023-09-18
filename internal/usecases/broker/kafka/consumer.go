package kafka

import (
	"context"
	"fmt"
	"strings"

	"github.com/antsrp/fio_service/internal/domain"
	ibroker "github.com/antsrp/fio_service/internal/interfaces/broker"
	"github.com/antsrp/fio_service/internal/mapper"
	kfk "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Consumer[T any] struct {
	reader   *kfk.Reader
	logger   *zap.Logger
	ctx      context.Context
	entities chan T
}

func NewConsumer[T any](host string, port int, topic string, groupId string, logger *zap.Logger,
	entities chan T, ctx context.Context) (*Consumer[T], error) {
	brokers := strings.Split(fmt.Sprintf("%s:%d", host, port), ",")
	c := &Consumer[T]{
		reader: kfk.NewReader(kfk.ReaderConfig{
			Brokers:     brokers,
			Topic:       topic,
			GroupID:     groupId,
			MinBytes:    10e3,
			MaxBytes:    10e6,
			StartOffset: kfk.LastOffset,
		}),
		logger:   logger,
		ctx:      ctx,
		entities: entities,
	}
	logger.Info("kafka consumer created")
	return c, nil
}

func (c Consumer[T]) Consume() {
	for c.ctx.Err() == nil {
		m, err := c.reader.ReadMessage(c.ctx)
		if err != nil {
			c.logger.Sugar().Errorf("can't consume data: %v", err.Error())
			continue
		}
		val, err := mapper.FromJSON[T](m.Value)
		if err != nil {
			c.logger.Sugar().Errorf("can't map kafka message to person type: %v", err.Error())
		} else {
			c.entities <- *val
		}
	}
}

func (c Consumer[T]) Close() error {
	c.logger.Info("kafka consumer closing")
	err := c.reader.Close()
	if err != nil {
		c.logger.Sugar().Errorf("can't close broker reader: %v", err.Error())
	} else {
		c.logger.Info("kafka consumer closed")
	}
	return err
}

var _ ibroker.Consumer[domain.PersonCommon] = Consumer[domain.PersonCommon]{}
