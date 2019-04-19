package counter

import (
	"errors"
	"testing"

	"github.com/wtask-go/auracounter/internal/api"
)

type repository struct {
	failEnsureSettings bool
	failGet            bool
	failIncrease       bool
	failSetSettings    bool
}

func (r *repository) EnsureSettings(_ int, _ *Settings) error {
	if r.failEnsureSettings {
		return errors.New("repository.EnsureSettings() failed")
	}
	return nil
}

func (r *repository) Get(_ int) (int, error) {
	if r.failGet {
		return 0, errors.New("repository.Get() failed")
	}
	return 0, nil
}

func (r *repository) Increase(_ int) (int, error) {
	if r.failIncrease {
		return 0, errors.New("repository.Increase() failed")
	}
	return 0, nil
}

func (r *repository) SetSettings(_ int, _ *Settings) error {
	if r.failSetSettings {
		return errors.New("repository.SetSettings() failed")
	}
	return nil
}

func TestServiceBuilder(t *testing.T) {
	cases := []struct {
		signature      string
		buildService   func() (api.CyclicCounterService, error)
		mustSuccessful bool
	}{
		{
			"NewCyclicCounterService(-1, nil)",
			func() (api.CyclicCounterService, error) { return NewCyclicCounterService(-1, nil) },
			false,
		},
		{
			"NewCyclicCounterService(0, nil)",
			func() (api.CyclicCounterService, error) { return NewCyclicCounterService(0, nil) },
			false,
		},
		{
			"NewCyclicCounterService(1, nil)",
			func() (api.CyclicCounterService, error) { return NewCyclicCounterService(1, nil) },
			false,
		},
		{
			"NewCyclicCounterService(1, &repository{failEnsureSettings: true})",
			func() (api.CyclicCounterService, error) {
				return NewCyclicCounterService(1, &repository{failEnsureSettings: true})
			},
			false,
		},
		{
			"NewCyclicCounterService(-1, &repository{})",
			func() (api.CyclicCounterService, error) { return NewCyclicCounterService(-1, &repository{}) },
			false,
		},
		{
			"NewCyclicCounterService(0, &repository{})",
			func() (api.CyclicCounterService, error) { return NewCyclicCounterService(0, &repository{}) },
			false,
		},
		{
			"NewCyclicCounterService(1, &repository{})",
			func() (api.CyclicCounterService, error) { return NewCyclicCounterService(1, &repository{}) },
			true,
		},
		{
			"NewCyclicCounterService(1, &repository{}, WithDefaults(DefaultSettings()))",
			func() (api.CyclicCounterService, error) {
				return NewCyclicCounterService(1, &repository{}, WithDefaults(DefaultSettings()))
			},
			true,
		},
		{
			"NewCyclicCounterService(1, &repository{}, WithDefaults(&Settings{Increment: -1}))",
			func() (api.CyclicCounterService, error) {
				return NewCyclicCounterService(1, &repository{}, WithDefaults(&Settings{Increment: -1}))
			},
			false,
		},
		{
			"NewCyclicCounterService(1, &repository{}, WithDefaults(&Settings{Increment: 0}))",
			func() (api.CyclicCounterService, error) {
				return NewCyclicCounterService(1, &repository{}, WithDefaults(&Settings{Increment: 0}))
			},
			true,
		},
		{
			"NewCyclicCounterService(1, &repository{}, WithDefaults(&Settings{StartFrom: -1, Upper: 1}))",
			func() (api.CyclicCounterService, error) {
				return NewCyclicCounterService(
					1,
					&repository{},
					WithDefaults(&Settings{StartFrom: -1, Upper: 1}),
				)
			},
			false,
		},
		{
			"NewCyclicCounterService(1, &repository{}, WithDefaults(&Settings{Lower: -1, Upper: 1, Increment: 1}))",
			func() (api.CyclicCounterService, error) {
				return NewCyclicCounterService(
					1,
					&repository{},
					WithDefaults(&Settings{Lower: -1, Upper: 1, Increment: 1}),
				)
			},
			true,
		},
		{
			"NewCyclicCounterService(1, &repository{}, WithDefaults(&Settings{Lower: 1, Upper: -1, Increment: 1}))",
			func() (api.CyclicCounterService, error) {
				return NewCyclicCounterService(
					1,
					&repository{},
					WithDefaults(&Settings{Lower: 1, Upper: -1, Increment: 1}),
				)
			},
			false,
		},
		{
			"NewCyclicCounterService(1, &repository{}, WithDefaults(&Settings{Lower: 1, Upper: -1, Increment: 1}))",
			func() (api.CyclicCounterService, error) {
				return NewCyclicCounterService(
					1,
					&repository{},
					WithDefaults(&Settings{Lower: 1, Upper: 3, Increment: 3}),
				)
			},
			false,
		},
	}

	for _, c := range cases {
		s, err := c.buildService()
		if c.mustSuccessful {
			if err != nil {
				t.Errorf("%s was expected to be successful but error occurred: %s", c.signature, err)
			} else if s == nil {
				t.Errorf("%s returns nil without error", c.signature)
			}
		} else {
			if err == nil {
				t.Errorf("%s was expected to be failed, but not", c.signature)
			} else if s != nil {
				t.Errorf("%s returns error with not nil service", c.signature)
			}
		}
	}
}

func TestGetCounterValue(t *testing.T) {
	service, err := NewCyclicCounterService(1, &repository{}) // good working repo
	if err != nil {
		// duplicates cases in TestServiceBuilder
		t.Errorf("Unexpected error when building service with non-failed repository: %+v", err)
	}

	intResult, apiErr := service.GetCounterValue()
	if intResult == nil {
		t.Errorf("GetCounterValue(): unexpected nil as result")
	}
	if apiErr != nil {
		t.Errorf("GetCounterValue(): unexpected API error %q", apiErr.Expose())
	}

	service, err = NewCyclicCounterService(1, &repository{failGet: true})
	if err != nil {
		// duplicates cases in TestServiceBuilder
		t.Errorf("Unexpected error when building service: %+v", err)
	}
	intResult, apiErr = service.GetCounterValue()
	if intResult != nil {
		t.Errorf("GetCounterValue(): unexpected non-nil result %+v", intResult)
	}
	if apiErr == nil {
		t.Errorf("GetCounterValue(): expected API error, but got nil")
	}
}

func TestIncreaseCounter(t *testing.T) {
	service, err := NewCyclicCounterService(1, &repository{}) // good working repo
	if err != nil {
		// duplicates cases in TestServiceBuilder
		t.Errorf("Unexpected error when building service with non-failed repository: %+v", err)
	}

	intResult, apiErr := service.IncreaseCounter()
	if intResult == nil {
		t.Errorf("IncreaseCounter(): unexpected nil as result")
	}
	if apiErr != nil {
		t.Errorf("IncreaseCounter(): unexpected API error %q", apiErr.Expose())
	}

	service, err = NewCyclicCounterService(1, &repository{failIncrease: true})
	if err != nil {
		// duplicates cases in TestServiceBuilder
		t.Errorf("Unexpected error when building service: %+v", err)
	}
	intResult, apiErr = service.IncreaseCounter()
	if intResult != nil {
		t.Errorf("IncreaseCounter(): unexpected non-nil result %+v", intResult)
	}
	if apiErr == nil {
		t.Errorf("IncreaseCounter(): expected API error, but got nil")
	}
}

func TestSetSettings(t *testing.T) {
	cases := []struct {
		signature          string
		setCounterSettings func(s api.CyclicCounterService) (*api.OKResult, *api.Error)
		mustSuccessful     bool // if true faulty must fail
	}{
		{
			// useless or sleepping counter
			"SetCounterSettings(0, 0, 0)",
			func(s api.CyclicCounterService) (*api.OKResult, *api.Error) {
				return s.SetCounterSettings(0, 0, 0)
			},
			true,
		},
		{
			// negative increment
			"SetCounterSettings(-1, 0, 1)",
			func(s api.CyclicCounterService) (*api.OKResult, *api.Error) {
				return s.SetCounterSettings(-1, 0, 1)
			},
			false,
		},
		{
			"SetCounterSettings(0, 0, 1)",
			func(s api.CyclicCounterService) (*api.OKResult, *api.Error) {
				return s.SetCounterSettings(0, 0, 1)
			},
			true,
		},
		{
			"SetCounterSettings(1, 0, 1)",
			func(s api.CyclicCounterService) (*api.OKResult, *api.Error) {
				return s.SetCounterSettings(1, 0, 1)
			},
			true,
		},
		{
			// increment is wider than range
			"SetCounterSettings(2, 0, 1)",
			func(s api.CyclicCounterService) (*api.OKResult, *api.Error) {
				return s.SetCounterSettings(2, 0, 1)
			},
			false,
		},
		{
			// invalid lower-upper range [2:0]
			"SetCounterSettings(2, 2, 0)",
			func(s api.CyclicCounterService) (*api.OKResult, *api.Error) {
				return s.SetCounterSettings(2, 2, 0)
			},
			false,
		},
	}

	good, err := NewCyclicCounterService(1, &repository{})
	if err != nil {
		// duplicates cases in TestServiceBuilder
		t.Errorf("Unexpected error when building service with non-failed repository: %+v", err)
	}
	faulty, err := NewCyclicCounterService(1, &repository{failSetSettings: true})
	if err != nil {
		// duplicates cases in TestServiceBuilder
		t.Errorf("Unexpected error when building service with faulty repository: %+v", err)
	}

	for _, c := range cases {
		okResult, err := c.setCounterSettings(good)
		faultyResult, faultyErr := c.setCounterSettings(faulty)
		if c.mustSuccessful {
			if okResult == nil {
				t.Errorf("%s: unexpected nil result for non-failed service", c.signature)
			}
			if err != nil {
				t.Errorf("%s: unexpected non-nil API error for non-failed service: %+v", c.signature, err)
			}
		} else {
			if okResult != nil {
				t.Errorf("%s: expected nil result for non-failed service, got: %+v ", c.signature, okResult)
			}
			if err == nil {
				t.Errorf("%s: unexpected nil API error for non-failed service", c.signature)
			}
		}
		// faulty service must fail in both cases:
		// - the same fail where good service fails
		// - fail, when accessing repository
		if faultyResult != nil {
			t.Errorf("%s: unexpected non-nil result for faulty service: %+v", c.signature, faultyResult)
		}
		if faultyErr == nil {
			t.Errorf("%s: unexpected nil API error for faulty service", c.signature)
		}
	}
}
