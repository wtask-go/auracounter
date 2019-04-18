package api

import (
	"github.com/pkg/errors"
)

type Error struct {
	Code     int // reserved
	Message  string
	Internal error // if nil it is not an internal error
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.Message == "" {
		if e.Internal != nil {
			return "unspecified internal error"
		}
		return "unspecified external error"
	}
	return e.Message
}

func (e *Error) Expose() error {
	if e == nil {
		return nil
	}
	if e.IsInternal() {
		return errors.WithMessage(e.Internal, e.Message)
	}
	return e
}

func (e *Error) IsInternal() bool {
	return e != nil && e.Internal != nil
}
