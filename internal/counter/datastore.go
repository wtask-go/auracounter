package counter

type Repository interface {
	EnsureSettings(counterID int, defaults *Settings) error
	Get(counterID int) (int, error)
	Increase(counterID int) (int, error)
	SetSettings(counterID, increment, upperLimit int) error
}

type Storage interface {
	EnsureLatestVersion() error
	Repository() Repository
	Close() error
}
