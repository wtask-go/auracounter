package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/wtask-go/auracounter/internal/httpcore/rest"

	"github.com/wtask-go/auracounter/internal/httpcore"
)

func main() {

	cfg, _ := configureFromCLI()
	// if err != nil {
	// 	os.Exit(1)
	// }

	l := log.New(os.Stdout, "auraserver ", log.LUTC|log.Ldate|log.Lmicroseconds|log.Ltime)
	l.Println("Hello!")

	server := newServer(cfg)

	shutdown, err := httpcore.LaunchServer(server, 3*time.Second)
	if err != nil {
		l.Println("ERR:", err)
		os.Exit(1)
	}
	l.Println("INFO:", "server successfully started")

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	<-sig
	if err := shutdown(10 * time.Second); err != nil {
		l.Println("ERR:", err)
	}
	l.Println("Bye!")
}

func newServer(cfg *Config) *http.Server {
	// nextRequestID := func() string {
	// 	return fmt.Sprintf("%d", time.Now().UnixNano())
	// }

	router := http.NewServeMux()
	router.Handle("/", index())

	return &http.Server{
		Addr: fmt.Sprintf("%s:%d", cfg.ServerAddress, cfg.ServerPort),
		// Handler: tracing(nextRequestID)(router),
		Handler: rest.NewCounterHandler(nil),
		// ErrorLog:     l,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
}

func index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Hello, World!")
	})
}

type key int

const (
	requestIDKey key = 0
)

func tracing(nextRequestID func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = nextRequestID()
			}
			ctx := context.WithValue(r.Context(), requestIDKey, requestID)
			w.Header().Set("X-Request-Id", requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
