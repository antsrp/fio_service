package broker

type Consumer[T any] interface {
	Consume()
	Close() error
}
