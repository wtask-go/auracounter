package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/wtask-go/auracounter/pkg/logging"

	"github.com/wtask-go/auracounter/internal/httpcore/rest"

	"github.com/wtask-go/auracounter/internal/api"

	"github.com/wtask-go/auracounter/internal/config"

	"github.com/wtask-go/auracounter/internal/counter"
	"github.com/wtask-go/auracounter/internal/counter/datastore/mysql"

	"github.com/wtask-go/auracounter/internal/httpcore"
)

func main() {
	logger := logging.NewStdout(logging.WithDecoration(logging.DefaultVerbosity("aurasrv", nil)))
	defer logger.Close()

	logger.Infof("Server is starting ...")

	storage := storageFactory(conf)
	defer storage.Close()
	server := newRESTServer(conf, counter.NewCounterService(storage.Repository()))

	shutdown, err := httpcore.LaunchServer(server, 3*time.Second)
	if err != nil {
		logger.Errorf("Can't launch server: %v", err)
		os.Exit(1)
	}
	logger.Infof("Server is ready!")

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	<-sig
	if err := shutdown(10 * time.Second); err != nil {
		logger.Errorf("Server shutdown failed: %v", err)
	}
	logger.Infof("Server has stopped, bye ( ᴗ_ ᴗ)")
}

func storageFactory(cfg *config.Application) counter.Storage {
	r, err := mysql.NewStorage(
		mysql.WithDSN(cfg.CounterDB.DSN()),
		mysql.WithCounterID(cfg.CounterID),
		mysql.WithTablePrefix(cfg.CounterDB.TablePrefix),
	)
	if err != nil {
		panic(err)
	}
	return r
}

func newRESTServer(cfg *config.Application, service api.CounterService) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.CounterREST.Host, cfg.CounterREST.Port),
		Handler: rest.NewCounterHandler(cfg.CounterREST.BaseURI, service),
		// ErrorLog:     l,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
}
