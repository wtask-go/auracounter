package api

// CyclicCounterService - represents interface for manage cyclic incremental counter.
type CyclicCounterService interface {
	// GetCounterValue - get current counter value
	GetCounterValue() (*IntValueResult, *Error)
	// IncreaseCounter - increase counter by increment, which set with settings and return new counter value.
	IncreaseCounter() (*IntValueResult, *Error)
	// SetCounterSettings - set the new settings for counter atomically
	SetCounterSettings(increment, lower, upper int) (*OKResult, *Error)
}

// IntValueResult - struct to return int value
type IntValueResult struct {
	Value int `json:"value"`
}

// OKResult - struct to return bool value (flag of success)
type OKResult struct {
	OK bool `json:"ok"`
}
