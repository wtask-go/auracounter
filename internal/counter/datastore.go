package counter

// Repository - contains methods to operate with counter
type Repository interface {
	// EnsureSettings - make sure settings are persisted for the counter with given ID.
	// If not, method will save default settings for counter.
	EnsureSettings(counterID int, defaults *Settings) error
	// Get - return current counter value.
	GetValue(counterID int) (int, error)
	// Increase - increase counter with increment which defined by settings.
	Increase(counterID int) (int, error)
	// SetSettings - set new counter settings
	SetSettings(counterID int, settings *Settings) error
}

// Storage - counter datastorage
type Storage interface {
	// EnsureLatest - make sure underlying database has latest version and is up-to-date to store counter.
	EnsureLatest() error
	// Repository - allows to explicitly expose the storage as a repository.
	Repository() Repository
	// Close - must close and free all used connections and resources.
	Close() error
}
