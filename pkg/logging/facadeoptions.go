package logging

import (
	"github.com/pkg/errors"
)

// WithDecoration - set decorator to format log messages.
func WithDecoration(d MessageDecorator) facadeOption {
	if d == nil {
		panic(errors.New("logging: can not use nil as MessageDecorator"))
	}
	return func(l *facade) {
		l.decorator = d
	}
}
