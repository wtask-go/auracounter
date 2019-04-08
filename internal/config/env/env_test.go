package env

import (
	"reflect"
	"strings"
	"github.com/wtask-go/auracounter/internal/config"
	"testing"
	"os"
	"github.com/joho/godotenv"
)

// up - reads env-file to collect var-names
// and loads it as envirinment
func loadEnv(file string, t *testing.T) ([]string, error) {
	env := "testdata/"+file
	m, err := godotenv.Read(env)
	if err != nil {
		return nil, err
	}
	if err = godotenv.Load(env); err != nil {
		return nil, err
	}
	vars := make([]string, 0, len(m))
	for v := range m {
		if v != "" {
			vars = append(vars, v)
		}
	}
	t.Logf("loaded environment (%d) from %q ", len(vars), env)
	return vars, nil
}

func cleanEnv(env []string, t *testing.T) {
	for _, name := range env {
		os.Unsetenv(name)
	} 
	t.Logf("environment cleaned (%d)", len(env))
}

// TestNewApplicationConfiguration - checks only loading config from environment vars,
// without final validation.
func TestNewApplicationConfig(t *testing.T) {
	cases := []struct{
		file string
		prefix string
		errMsg string
		expConfig *config.Application
	} {
		{
			// may need to remove this test due probably it can overwrite real environment
			"correct.env", 
			"",
			"",
			&config.Application{
				CounterREST: config.HTTPServer{
					Host: "", 
					Port: 33333,
					BaseURI: "/counter/v1/",
				},
				CounterDB: config.Database{
					Type:        "mysql",
					Host:        "127.0.0.1",
					Port:        3306,
					Name:        "database",
					User:        "user",
					Password:    "password",
					Options:     "parseTime=true&timeout=3m",
					TablePrefix: "",
				},
				CounterID: 1,
			},
		},
		{
			"correct-with-prefix.env", 
			"ENVTEST_",
			"",
			&config.Application{
				CounterREST: config.HTTPServer{
					Host: "", 
					Port: 33333,
					BaseURI: "/counter/v1/",
				},
				CounterDB: config.Database{
					Type:        "mysql",
					Host:        "127.0.0.1",
					Port:        3306,
					Name:        "database",
					User:        "user",
					Password:    "password",
					Options:     "parseTime=true&timeout=3m",
					TablePrefix: "",
				},
				CounterID: 1,
			},
		},
		{
			"incorrect-due-prefix.env", "ENVTEST_", "error: \"ENVTEST_COUNTER_ID\" is required (int)", nil,
		},
		{
			"incorrect-due-rest-port.env", "", "error: optional \"COUNTER_REST_PORT\" is expected as int", nil,
		},
		{
			"incorrect-due-db-host.env", "", "error: \"COUNTER_DB_HOST\" is required (string)", nil,
		},
		{
			"incorrect-due-db-port.env", "", "error: required \"COUNTER_DB_PORT\" is expected as int", nil,
		},
		{
			"incorrect-due-db-name.env", "", "error: \"COUNTER_DB_NAME\" is required (string)", nil,
		},
		{
			"incorrect-due-db-user.env", "", "error: \"COUNTER_DB_USER\" is required (string)", nil,
		},
		{
			"incorrect-due-db-password.env", "", "error: \"COUNTER_DB_PASSWORD\" is required (string)", nil,
		},
		{
			"incorrect-due-counter-id-1.env", "", "error: \"COUNTER_ID\" is required (int)", nil,
		},
		{
			"incorrect-due-counter-id-2.env", "", "error: required \"COUNTER_ID\" is expected as int", nil,
		},
	}

	for _, c := range cases {
		envVars, err := loadEnv(c.file, t)
		if err != nil {
			t.Fatalf("unable load %q", "testdata/"+c.file)
		}

		config, err := NewApplicationConfig(c.prefix)
		if err != nil {
			if c.errMsg == "" {
				t.Errorf("Unexpected error: %q", err)
			} else if strings.Index(err.Error(), c.errMsg) == -1{
				t.Errorf("Expected error message: %q is not contained in: %q", c.errMsg, err)
			}
		} else if c.errMsg != "" {
			t.Errorf("Expected error will contain message: %q, but nothing happened", c.errMsg)
		}

		if err == nil && !reflect.DeepEqual(c.expConfig, config) {
			t.Errorf("Expected config: %+v, got: %+v", c.expConfig, config)
		}

		cleanEnv(envVars, t)
	}
}