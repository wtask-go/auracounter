package api

// ErrorInvalidArgument - specific error to help make right response
type ErrorInvalidArgument struct {
	Message string
}

func (e *ErrorInvalidArgument) Error() string {
	if e.Message == "" {
		return "invalid argument value"
	}
	return e.Message
}

// IsRequestError - checks that the error in the incoming request data
func IsRequestError(e error) bool {
	switch e.(type) {
	case *ErrorInvalidArgument:
		return true
	}
	return false
}
