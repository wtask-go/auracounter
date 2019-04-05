package counter

type Repository interface {
	GetNumber() (int, error)
	IncrementNumber() (int, error)
	SetSettings(delta, max int) error
}

type Storage interface {
	Repository() Repository
	Close() error
}
