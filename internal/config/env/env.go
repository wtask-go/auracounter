package env

import (
	"fmt"
	"os"
	"strconv"

	"github.com/wtask-go/auracounter/internal/config"

	"github.com/pkg/errors"
)

// NewApplicationConfig - loads application configuration from system environment.
// Parameter `prefix` defines common prefix for all environment variables names.
// See config examples in the provided env-files.
func NewApplicationConfig(prefix string) (cfg *config.Application, err error) {
	defer func() {
		// need to return error with expected environment var name
		if e := recover(); e != nil {
			err, _ = e.(error)
			if err != nil {
				err = errors.Wrapf(err, "env.NewApplicationConfig error")
			} else {
				err = errors.Errorf("env.NewApplicationConfig failed: %v", e)
			}
			cfg = nil
		}
	}()
	// p - return complete name of var
	p := func(name string) string {
		return prefix + "COUNTER_" + name
	}
	return &config.Application{
		CounterREST: config.HTTPServer{
			Host: optionalString(p("REST_HOST"), ""),
			Port: optionalInt(p("REST_PORT"), 33333),
		},
		CounterDB: config.Database{
			Type:        "mysql",
			Host:        requiredString(p("DB_HOST")),
			Port:        requiredInt(p("DB_PORT")),
			Name:        requiredString(p("DB_NAME")),
			User:        requiredString(p("DB_USER")),
			Password:    requiredString(p("DB_PASSWORD")),
			Options:     optionalString(p("DB_OPTIONS"), "parseTime=true"),
			TablePrefix: optionalString(p("DB_TABLE_PREFIX"), ""),
		},
		CounterID: requiredInt(p("ID")),
	}, nil
}

// optionalString - obtain string value from environment.
// If var is not defined returns defaults.
func optionalString(varname, defaults string) string {
	val, ok := os.LookupEnv(varname)
	if !ok {
		return defaults
	}
	return val
}

// lookupString - obtain string value from environment.
// Panics, if var is not defined.
func requiredString(varname string) string {
	val, ok := os.LookupEnv(varname)
	if !ok {
		panic(fmt.Errorf("%q is required (string)", varname))
	}
	return val
}

// optionalInt - obtain integer value from environment.
// Panics, if var defined, but can not be converted into int.
func optionalInt(varname string, defaults int) int {
	str, ok := os.LookupEnv(varname)
	if !ok {
		return defaults
	}
	val, err := strconv.Atoi(str)
	if err != nil {
		panic(errors.Wrapf(err, "optional %q is expected as int", varname))
	}
	return val
}

// requiredInt - obtain integer value from environment.
// Panics, if var is not defined or can not be converted into int.
func requiredInt(varname string) int {
	str, ok := os.LookupEnv(varname)
	if !ok {
		panic(fmt.Errorf("%q is required (int)", varname))
	}
	val, err := strconv.Atoi(str)
	if err != nil {
		panic(errors.Wrapf(err, "required %q is expected as int", varname))
	}
	return val
}
