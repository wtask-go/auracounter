package logging

import (
	"sync"

	"github.com/pkg/errors"
)

// Chain represents simple logging chain.
// It is support null/empty/blackhole loggin out of the box.
// You can use Chain{} or logging.Chain{nil} for it.
type Chain []Interface

// eachFacade - call visit() func for every chained logging facade.
func (c Chain) eachFacade(visit func(f Facade)) {
	wg := sync.WaitGroup{}
	for _, facade := range c {
		if facade == nil {
			continue
		}
		wg.Add(1)
		go func(f Facade) {
			defer wg.Done()
			visit(f)
		}(facade)
	}
	wg.Wait()
}

// Infof - logging.Facade implementation.
// Visits chained items and calls its Infof() method.
func (c Chain) Infof(format string, a ...interface{}) {
	c.eachFacade(func(f Facade) { f.Infof(format, a...) })
}

// Errorfn - logging.Facade implementation.
// Visits chained items and calls its Errorf() method.
func (c Chain) Errorf(format string, a ...interface{}) {
	c.eachFacade(func(f Facade) { f.Errorf(format, a...) })
}

// Close - closes all logging interfaces in the chain
func (c Chain) Close() error {
	errs := []error{}
	for _, logger := range c {
		if logger == nil {
			continue
		}
		if err := errors.WithMessagef(logger.Close(), "%T", logger); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.Errorf("logging.Chain closed with errors: %v", errs)
	}
	return nil
}
