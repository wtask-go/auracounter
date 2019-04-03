package logging

import (
	"testing"
)

func TestEmptyChain(t *testing.T) {
	// test will fail on underlying panic

	// empty chain works as blackhole (null logger)
	blackhole := Chain{}
	blackhole.Errorf("Some error")
	blackhole.Infof("Some information")
	blackhole.Tracef("Some debug or trace point")
	if err := blackhole.Close(); err != nil {
		t.Errorf("Closing empty chain failed: %+v", err)
	}

	// literal contains redundant nils, but chain must also work
	blackhole = Chain{nil, nil, nil}
	blackhole.Errorf("Some error")
	blackhole.Infof("Some information")
	blackhole.Tracef("Some debug or trace point")
	if err := blackhole.Close(); err != nil {
		t.Errorf("Closing redundant empty chain failed: %+v", err)
	}
}
