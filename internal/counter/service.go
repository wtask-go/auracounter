package counter

import (
	"github.com/pkg/errors"
	"github.com/wtask-go/auracounter/internal/api"
)

type serviceOption func(*service)

// apply - apply given options for stream.
func (s *service) apply(options ...serviceOption) *service {
	if s == nil {
		return nil
	}
	for _, o := range options {
		if o != nil {
			o(s)
		}
	}
	return s
}

func WithDefaults(defaults *Settings) serviceOption {
	if defaults == nil {
		panic(errors.New("counter.WithDefaults: unable to use nil settings"))
	}
	if err := defaults.verify(); err != nil {
		panic(err)
	}
	return func(s *service) {
		s.defaults = defaults
	}
}

// NewCounterService - create new instance of api.CounterService implementation.
func NewCounterService(counterID int, r Repository, options ...serviceOption) api.CounterService {
	// will reserve zero ID
	if counterID < 0 {
		panic(errors.Errorf("counter.NewCounterService: invalid counter ID (%d)", counterID))
	}
	if r == nil {
		panic(errors.New("counter.NewCounterService: unable to use nil Repository"))
	}
	s := (&service{
		repo:      r,
		counterID: counterID,
		defaults:  DefaultSettings(),
	}).apply(options...)
	if err := s.repo.EnsureSettings(s.counterID, s.defaults);err!=nil {
		panic(errors.Errorf("counter.NewCounterService: unable to use nil Repository"))
	}
	return s
}

type service struct {
	repo      Repository
	counterID int
	defaults  *Settings
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
