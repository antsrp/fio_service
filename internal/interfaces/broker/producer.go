package broker

type Producer[T any] interface {
	Produce(key string, entity T) error
	Close() error
}
