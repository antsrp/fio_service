package kafka

import (
	"context"
	"fmt"

	"github.com/antsrp/fio_service/internal/domain"
	ibroker "github.com/antsrp/fio_service/internal/interfaces/broker"
	"github.com/antsrp/fio_service/internal/mapper"
	kfk "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Producer[T any] struct {
	writer *kfk.Writer
	logger *zap.Logger
}

func NewProducer[T any](host string, port int, topic string, logger *zap.Logger) (*Producer[T], error) {
	url := fmt.Sprintf("%s:%d", host, port)
	p := &Producer[T]{
		writer: &kfk.Writer{
			Addr:     kfk.TCP(url),
			Topic:    topic,
			Balancer: &kfk.LeastBytes{},
		},
		logger: logger,
	}
	if err := p.writer.WriteMessages(context.Background()); err != nil {
		return nil, err
	}
	logger.Info("kafka producer created")

	return p, nil
}

func (p Producer[T]) Produce(key string, entity T) error {
	data, err := mapper.ToJSON(entity, &mapper.Indent{})
	if err != nil {
		return fmt.Errorf("can't marshal data: %w", err)
	}
	msg := kfk.Message{
		Key:   []byte(key),
		Value: data,
	}

	if err := p.writer.WriteMessages(context.Background(), msg); err != nil {
		return fmt.Errorf("can't write message into broker's topic %s: %w", p.writer.Topic, err)
	}
	return nil
}

func (p Producer[T]) Close() error {
	p.logger.Info("kafka producer closing")
	err := p.writer.Close()
	if err != nil {
		p.logger.Sugar().Errorf("can't close broker writer: %v", err.Error())
	} else {
		p.logger.Info("kafka producer closed")
	}
	return err
}

var _ ibroker.Producer[domain.PersonCommon] = Producer[domain.PersonCommon]{}
