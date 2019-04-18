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

// Error - go Error interface implementation.
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

// Expose - expose Error struct as solid standard error.
func (e *Error) Expose() error {
	if e == nil {
		return nil
	}
	if e.IsInternal() {
		return errors.WithMessage(e.Internal, e.Message)
	}
	return e
}

// IsInternal - checks if API Error is internal (server/infrastructure) error (if true) or client error (if false).
func (e *Error) IsInternal() bool {
	return e != nil && e.Internal != nil
}
