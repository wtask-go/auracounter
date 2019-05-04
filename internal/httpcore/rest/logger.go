package rest

import (
	"fmt"
	"net/http"
)

// Logger - interface used by rest-package to log two types of messages.
type Logger interface {
	Error(a ...interface{})
	Info(a ...interface{})
}

// logInfo - helps to log info messages.
func logInfo(l Logger, a ...interface{}) {
	if l == nil {
		return
	}
	l.Info(a...)
}

// logError - helps to log error messages.
func logError(l Logger, a ...interface{}) {
	if l == nil {
		return
	}
	l.Error(a...)
}

// formatRequest - formats request attributes as solid string.
func formatRequest(r *http.Request) string {
	return fmt.Sprintf("%s %s %s %s %s", r.Proto, r.Method, r.URL, r.RemoteAddr, r.UserAgent())
}

// formatError - formats the error with +v specifier and returns result in quotes.
func formatError(e error) string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("%q", fmt.Sprintf("%+v", e))
}
