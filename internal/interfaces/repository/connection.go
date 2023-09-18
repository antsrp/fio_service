package repository

type DBConnection interface {
	Check() error
	Close() error
}
