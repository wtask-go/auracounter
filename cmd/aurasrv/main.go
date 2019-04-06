package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/wtask-go/auracounter/internal/logging"

	"github.com/wtask-go/auracounter/internal/counter"
	"github.com/wtask-go/auracounter/internal/counter/datastore/mysql"

	"github.com/wtask-go/auracounter/internal/httpcore/rest"

	"github.com/wtask-go/auracounter/internal/httpcore"
)

func main() {

	cfg, err := configureFromCLI()
	if err != nil {
		fmt.Println("STARTUP ERROR:", err)
		os.Exit(1)
	}

	logger := logging.NewStdOut(logging.WithPrefix("aurasrv "), logging.WithTrace(false))
	defer logger.Close()

	logger.Infof("Server is starting ...")
	logger.Errorf("Check log has %s tag", "ERR")

	// l := log.New(os.Stdout, "auraserver ", log.LUTC|log.Ldate|log.Lmicroseconds|log.Ltime)
	l := logger.ExposeLogger(" ")
	l.Println("Hello!")

	storage := storageFactory(cfg)
	defer storage.Close()
	service := counter.NewCounterService(storage.Repository())
	handleREST := rest.NewCounterHandler(service)
	server := newServer(cfg.ServerAddress, cfg.ServerPort, handleREST)

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

func storageFactory(cfg *Config) counter.Storage {
	r, err := mysql.NewStorage(
		mysql.WithDSN(cfg.StorageDSN),
		mysql.WithCounterID(cfg.CounterID),
		mysql.WithTablePrefix("aura_"),
	)
	if err != nil {
		panic(err)
	}
	return r
}

func newServer(addr string, port int, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf("%s:%d", addr, port),
		Handler: handler,
		// ErrorLog:     l,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
}
