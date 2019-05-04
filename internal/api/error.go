package api

import (
	"github.com/pkg/errors"
)

// Error - API error representation.
type Error struct {
	// Code - error code, reserved
	Code int // reserved
	// Message - public error message (without infrastructure details)
	Message string
	// Internal - complete internal error if it was
	Internal error // if nil it is not an internal error
}

// Error - returns error message without references to internal details if any.
func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.Message == "" {
		if e.Internal != nil {
			return "Internal error"
		}
		return "Unspecified error"
	}
	return e.Message
}

// ExposeError - complements API error with internal details and returns as standard error interface.
// Should not be used for client-side errors.
func (e *Error) ExposeError() error {
	if e == nil {
		return nil
	}
	if e.IsInternal() {
		return errors.WithMessage(e.Internal, e.Error())
	}
	return e
}

// IsInternal - checks if API Error is internal (server/infrastructure) or client error.
func (e *Error) IsInternal() bool {
	return e != nil && e.Internal != nil
}
