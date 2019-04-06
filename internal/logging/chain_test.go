package logging

import (
	"testing"
)

func TestEmptyChain(t *testing.T) {
	// test will fail on underlying panic

	// empty chain works as blackhole (null logger)
	blackhole := Chain{}
	blackhole.Errorfn("Some error")
	blackhole.Infofn("Some information")
	if err := blackhole.Close(); err != nil {
		t.Errorf("Closing empty chain failed: %+v", err)
	}

	// literal contains redundant nils, but chain must also work
	blackhole = Chain{nil, nil, nil}
	blackhole.Errorfn("Some error")
	blackhole.Infofn("Some information")
	if err := blackhole.Close(); err != nil {
		t.Errorf("Closing redundant empty chain failed: %+v", err)
	}
}
