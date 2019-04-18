package counter

import (
	"github.com/pkg/errors"
	"github.com/wtask-go/auracounter/internal/api"
)

type (
	// service - struct to implement api.CyclicCounterService interface
	service struct {
		repo      Repository
		counterID int
		defaults  *Settings
	}

	// serviceOption - func to set unexported optional field of service struct
	serviceOption func(*service) error
)

// invalidOption - helper to return error from option builder
func invalidOption(e error) serviceOption {
	return func(*service) error {
		return e
	}
}

// setup - applies options, but stops after first error
func (s *service) setup(options ...serviceOption) error {
	if s == nil {
		return nil
	}
	for _, setter := range options {
		if setter == nil {
			continue
		}
		if err := setter(s); err != nil {
			return err
		}
	}
	return nil
}

// WithDefaults - sets default service settings.
// Service will use given settings if there are not persisted before.
func WithDefaults(settings *Settings) serviceOption {
	if settings == nil {
		return invalidOption(errors.New("counter.WithDefaults: unable to use nil settings"))
	}
	if err := settings.verify(); err != nil {
		return invalidOption(errors.WithMessage(err, "counter.WithDefaults: invalid settings"))
	}
	return func(s *service) error {
		s.defaults = settings
		return nil
	}
}

// NewCyclicCounterService - create new instance of api.CyclicCounterService implementation.
func NewCyclicCounterService(counterID int, r Repository, options ...serviceOption) (api.CyclicCounterService, error) {
	// required params
	if counterID <= 0 {
		return nil, errors.Errorf("counter.NewCyclicCounterService: invalid counter ID (%d)", counterID)
	}
	if r == nil {
		return nil, errors.New("counter.NewCyclicCounterService: unable to use nil as Repository")
	}
	s := &service{
		repo:      r,
		counterID: counterID,
		defaults:  DefaultSettings(),
	}
	// options
	if err := s.setup(options...); err != nil {
		return nil, errors.WithMessage(err, "counter.NewCyclicCounterService: setup error")
	}
	if err := s.repo.EnsureSettings(s.counterID, s.defaults); err != nil {
		return nil,
			errors.WithMessage(err, "counter.NewCyclicCounterService: unable to ensure persisted counter settings")
	}
	return s, nil
}

func (s *service) GetCounterValue() (*api.IntValueResult, *api.Error) {
	value, err := s.repo.Get(s.counterID)
	if err != nil {
		// TODO log internal error
		return nil, &api.Error{Message: "failed to get counter value", Internal: err}
	}
	return &api.IntValueResult{Value: value}, nil
}

func (s *service) IncreaseCounter() (*api.IntValueResult, *api.Error) {
	value, err := s.repo.Increase(s.counterID)
	if err != nil {
		// TODO log internal error
		return nil, &api.Error{Message: "failed to increase counter", Internal: err}
	}
	return &api.IntValueResult{Value: value}, nil
}

func (s *service) SetCounterSettings(increment, lower, upper int) (*api.OKResult, *api.Error) {
	err := (&Settings{
		StartFrom: lower, // we disallow to set start in this version
		Increment: increment,
		Lower:     lower, // for API v1 expected 0 always
		Upper:     upper,
	}).verify()
	if err != nil {
		return nil, &api.Error{Message: err.Error()}
	}
	if err := s.repo.SetSettings(increment, lower, upper); err != nil {
		return nil, &api.Error{Message: "failed to set new settings", Internal: err}
	}
	return &api.OKResult{OK: true}, nil
}
