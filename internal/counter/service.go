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

	serviceOption func() (func(*service), error)
)

// failedOption - helper to expose error from option builder
func failedOption(err error) serviceOption {
	return func() (func(*service), error) {
		return nil, err
	}
}

// properOption - helper to expose setter from option builder
func properOption(setter func(*service)) serviceOption {
	return func() (func(*service), error) {
		return setter, nil
	}
}

// setup - applies options, but stops on first error
func (s *service) setup(options ...serviceOption) error {
	if s == nil {
		return nil
	}
	for _, option := range options {
		if option == nil {
			continue
		}
		setter, err := option()
		if err != nil {
			return err
		}
		if setter != nil {
			setter(s)
		}
	}
	return nil
}

// WithDefaults - sets default service settings.
// Service will use given settings if there are not persisted before.
func WithDefaults(settings *Settings) serviceOption {
	if settings == nil {
		return failedOption(errors.New("counter.WithDefaults: unable to use nil settings"))
	}
	if err := settings.verify(); err != nil {
		return failedOption(errors.WithMessage(err, "counter.WithDefaults: invalid settings"))
	}
	return properOption(func(s *service) {
		s.defaults = settings
	})
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
	value, err := s.repo.GetValue(s.counterID)
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
	settings := &Settings{
		StartFrom: lower, // we disallow to set start in this version
		Increment: increment,
		Lower:     lower, // for API v1 expected 0 always
		Upper:     upper,
	}
	if err := settings.verify(); err != nil {
		return nil, &api.Error{Message: err.Error()}
	}
	if err := s.repo.SetSettings(s.counterID, settings); err != nil {
		return nil, &api.Error{Message: "failed to set new settings", Internal: err}
	}
	return &api.OKResult{OK: true}, nil
}
