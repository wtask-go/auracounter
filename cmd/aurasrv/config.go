package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// Config for auraserver
type Config struct {
	ServerAddress string
	ServerPort    int
	StorageDSN    string
	CounterID     int
}

// verifyConfig - checks config for errors
func verifyConfig(cfg *Config) error {
	if strings.Index(cfg.ServerAddress, ":") != -1 {
		return errors.Errorf("invalid address, must not contain %q", ":")
	}
	if cfg.ServerPort < 0 {
		return errors.Errorf("invalid server port %d", cfg.ServerPort)
	}
	if cfg.StorageDSN == "" {
		return errors.New("empty DSN")
	}
	if !strings.HasPrefix(cfg.StorageDSN, "mysql://") {
		return errors.New("invalid DSN, must prefixed with mysql://")
	}
	if cfg.CounterID < 1 {
		return errors.Errorf("invalid counter ID %d", cfg.CounterID)
	}
	return nil
}

func configureFromCLI() (*Config, error) {
	cfg := Config{}
	flag.StringVar(
		&cfg.ServerAddress,
		"addr",
		``,
		"Optional: server ip-address, for example `0.0.0.0` or host name.\n"+
			"If value was not set, server will bind all available addresses",
	)
	flag.IntVar(&cfg.ServerPort, "port", 33333, "Optional: server will listen on this port number")
	flag.StringVar(
		&cfg.StorageDSN,
		"dsn",
		"",
		"Required: DSN string to connect with DB. Only MySQL is supported.\n"+
			"Format: `mysql://user:pass@tcp(host-or-ip:port)/database`",
	)
	flag.IntVar(&cfg.CounterID, "cid", 1, "Required: Counter ID to maintain.\n")

	if !flag.Parsed() {
		flag.Parse()
	}
	err := verifyConfig(&cfg)
	if err != nil {
		fmt.Fprintf(flag.CommandLine.Output(), "Startup error: %s\n", err)
		flag.Usage()
	}
	return &cfg, err
}
