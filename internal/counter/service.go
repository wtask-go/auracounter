package counter

import (
	"github.com/pkg/errors"
	"github.com/wtask-go/auracounter/internal/api"
)

// NewCounterService - create new instance of api.CounterService implementation.
func NewCounterService(r Repository) api.CounterService {
	if r == nil {
		panic(errors.New("counter.NewCounterService: counter.Repository is not implemented"))
	}
	return &service{
		repo: r,
	}
}

type service struct {
	repo Repository
}

func (s *service) GetNumber() (*api.GetNumberResult, error) {
	num, err := s.repo.GetNumber()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get number from repository")
	}
	return &api.GetNumberResult{Value: num}, nil
}

func (s *service) IncrementNumber() (*api.IncrementNumberResult, error) {
	num, err := s.repo.IncrementNumber()
	if err != nil {
		return nil, errors.Wrap(err, "failed to increment number due to repository")
	}
	return &api.IncrementNumberResult{Value: num}, nil
}

func (s *service) SetSettings(delta, max int) (*api.SetSettingsResult, error) {
	if delta <= 0 {
		return nil, &api.ErrorInvalidArgument{"invalid counter delta"}
	}

	if max <= 0 {
		return nil, &api.ErrorInvalidArgument{"invalid max value for counter"}
	}

	if max < delta {
		return nil, &api.ErrorInvalidArgument{"mutually exclusive arguments"}
	}

	if err := s.repo.SetSettings(delta, max); err != nil {
		return nil, errors.Wrap(err, "failed set new settings due to repository")
	}
	return &api.SetSettingsResult{OK: true}, nil
}
