package main

import (
	"github.com/wtask-go/auracounter/internal/config"

	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/wtask-go/auracounter/internal/config/env"
)

const (
	envVarPrefix = "AURA_"
)

var (
	conf *config.Application
)

func init() {
	var err error
	usage := "aurasrv\nStarts REST HTTP server to maintain distributed counter.\n"
	envFile := ""
	help := false
	flag.StringVar(
		&envFile,
		"config",
		"",
		"Absolute path (optional) to application config in ENV-format (.env file).\n"+
			"With this parameter you can override config based on the system environment.\n",
	)
	flag.BoolVar(&help, "help", false, "Prints usage help.")
	flag.BoolVar(&help, "h", false, "")

	if !flag.Parsed() {
		flag.Parse()
	}
	out := flag.CommandLine.Output()
	if help {
		fmt.Fprintf(out, usage)
		flag.Usage()
		os.Exit(0)
	}
	if envFile != "" {
		if err = godotenv.Load(envFile); err != nil {
			fmt.Fprintf(out, "Can not load environment (%s): %s\n", envFile, err)
			os.Exit(1)
		}
	}

	if conf, err = env.NewApplicationConfig(envVarPrefix); err != nil {
		fmt.Fprintf(out, "Can not prepare config: %s\n", err)
		fmt.Fprintf(out, "Check usage with -help option. \n")
		os.Exit(1)
	}
}
