package httpcore

import (
	"context"
	"net/http"
	"time"
)

// StartServer - starts http server in background or return startup error.
func StartServer(server *http.Server, startupTimeout time.Duration) error {
	fail := make(chan error, 1)
	go func() {
		fail <- server.ListenAndServe()
		close(fail)
	}()
	select {
	case err := <-fail:
		return err
	case <-time.After(startupTimeout):
		return nil
	}
}

// StopServer - gracefully stops http server.
func StopServer(server *http.Server, shutdownTimeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	server.SetKeepAlivesEnabled(false)
	return server.Shutdown(ctx)
}

// LaunchServer - starts server in background.
// Parameter `timeout` is used as а period when startup errors are expected.
// First returned value is a shutdown function for started server.
// Second returned value is a startup error.
// Shutdown function will return error if server will stop with it.
// You must pass timeout into shutdown function so the server has time to stop.
func LaunchServer(server *http.Server, timeout time.Duration) (shutdown func(timeout time.Duration) error, startup error) {
	if err := StartServer(server, timeout); err != nil {
		return nil, err
	}
	quit := make(chan time.Duration)
	fail := make(chan error, 1)
	go func() {
		defer close(fail)
		d := <-quit
		fail <- StopServer(server, d)

	}()
	return func(d time.Duration) error {
		defer close(quit)
		quit <- d // send shutdown timeout
		return <-fail
	}, nil
}
