package api

type CyclicCounterService interface {
	GetCounterValue() (*IntValueResult, *Error)
	IncreaseCounter() (*IntValueResult, *Error)
	SetCounterSettings(increment, lower, upper int) (*OKResult, *Error)
}

type IntValueResult struct {
	Value int `json:"value"`
}

type OKResult struct {
	OK bool `json:"ok"`
}
