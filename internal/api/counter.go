package api

type IncrementalCounterService interface {
	GetCounterValue() (*GetCounterValue, error)
	IncreaseCounter() (*IncreaseCounterResult, error)
	SetCounterSettings(increment, upperLimit int) (*SetSettingsResult, error)
}

type GetCounterValue struct {
	Value int `json:"value"`
}

type IncreaseCounterResult struct {
	Value int `json:"value"`
}

type SetSettingsResult struct {
	OK bool `json:"ok"`
}
