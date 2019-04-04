package api

type CounterService interface {
	GetNumber() (*GetNumberResult, error)
	IncrementNumber() (*IncrementNumberResult, error)
	SetSettings(step, max int) (*SetSettingsResult, error)
}

type GetNumberResult struct {
	Value int `json:"value"`
}

type IncrementNumberResult struct {
	Value int `json:"value"`
}

type SetSettingsResult struct {
	OK bool `json:"ok"`
}
