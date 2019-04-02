package logging

import (
	"sync"
)

// Chain represents simple logging chain.
// It is interesting, chain pattern for logging supports null/empty/blackhole loggin out of the box.
// You can use empty chain for it, or logging.Chain{nil}
type Chain []Facade

// async - iterate chained items and start visit() func in go-routine.
func (c Chain) async(visit func(f Facade)) {
	wg := sync.WaitGroup{}
	for _, logger := range c {
		if logger == nil {
			continue
		}
		wg.Add(1)
		go func(f Facade) {
			defer wg.Done()
			visit(f)
		}(logger)
	}
	wg.Wait()
}

// Tracef - logging.Facade implementation.
// Visits chained items and calls its Tracef() method.
func (c Chain) Tracef(format string, a ...interface{}) {
	c.async(func(f Facade) { f.Tracef(format, a...) })
}

// Infof - logging.Facade implementation.
// Visits chained items and calls its Infof() method.
func (c Chain) Infof(format string, a ...interface{}) {
	c.async(func(f Facade) { f.Infof(format, a...) })
}

// Errorf - logging.Facade implementation.
// Visits chained items and calls its Errorf() method.
func (c Chain) Errorf(format string, a ...interface{}) {

}

func (c Chain) Close() error {
	return nil
}
