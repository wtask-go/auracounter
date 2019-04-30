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
	var exitCode int
	defer func() {
		// the latest deferred call after return from main
		os.Exit(exitCode)
	}()

	logger := logging.NewStdout(logging.WithDefaultDecoration("aurasrv", nil))
	defer logger.Close()

	logger.Infof("Server initialization started ...")

	storage, err := storageFactory(conf)
	if err != nil {
		logger.Errorf("Can't initialize storage: %v", err)
		exitCode = 1
		return
	}
	defer storage.Close()

	if err = storage.EnsureLatest(); err != nil {
		logger.Errorf("Can't ensure storage has latest version: %v", err)
		exitCode = 1
		return
	}

	service, err := counter.NewCyclicCounterService(conf.CounterID, storage.Repository())
	if err != nil {
		logger.Errorf("Can't initialize counter service: %v", err)
		exitCode = 1
		return
	}

	logger.Infof("Initialization done, server is starting ...")

	server := newRESTServer(conf, service)

	shutdown, err := httpcore.LaunchServer(server, 3*time.Second)
	if err != nil {
		logger.Errorf("Can't launch server: %v", err)
		exitCode = 1
		return
	}
	logger.Infof("Server is ready!")

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	<-sig
	if err := shutdown(10 * time.Second); err != nil {
		logger.Errorf("Server shutdown failed: %v", err)
		exitCode = 2
	}
	logger.Infof("Server has stopped, bye ( ᴗ_ ᴗ)")
}

func storageFactory(cfg *config.Application) (counter.Storage, error) {
	return mysql.NewStorage(cfg.CounterDB.DSN(), mysql.WithTablePrefix(cfg.CounterDB.TablePrefix))
}

func newRESTServer(cfg *config.Application, service api.CyclicCounterService) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.CounterREST.Host, cfg.CounterREST.Port),
		Handler: rest.NewCounterHandler(cfg.CounterREST.BaseURI, service),
		// ErrorLog:     l,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
}
