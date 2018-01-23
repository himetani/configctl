package workspace

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func setupForInit() error {
	testHome := os.Getenv("TEST_CONFIGCTL_HOME")
	if testHome == "" {
		return errors.New("Set TEST_CONFIGCTL_HOME environment variable to test")
	}
	os.Setenv("CONFIGCTL_HOME", testHome)
	return nil
}

func teardownForInit() {
	if err := os.RemoveAll(os.Getenv("CONFIGCTL_HOME")); err != nil {
		log.Printf("Unexpected error happened at teardown: %s\n", err)
	}
}

func TestInit(t *testing.T) {
	if err := setupForInit(); err != nil {
		t.Fatal(err)
	}

	configCtlHome := os.Getenv("CONFIGCTL_HOME")

	if err := Init(); err != nil {
		t.Errorf("Unexpected error happend: %s\n", err)
	}

	defer teardownForInit()

	if _, err := os.Stat(configCtlHome); os.IsNotExist(err) {
		t.Errorf("$CONFIGCTL_HOME dir is not created. %s\n", err)
	}

	if _, err := os.Stat(filepath.Join(configCtlHome, "configs")); os.IsNotExist(err) {
		t.Errorf("$CONFIGCTL_HOME/configs is not created. %s\n", err)
	}
}

func TestGetConfigs(t *testing.T) {
	configCtlHome = filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "himetani", "configctl", "workspace", "testdata")

	t.Run("GetConfigs", func(t *testing.T) {
		type data struct {
			TestName string
			Result   []string
		}

		tests := []data{
			{"Success", []string{"valid"}},
		}

		for i, test := range tests {
			result := GetConfigs()
			if !reflect.DeepEqual(result, test.Result) {
				t.Errorf("Test #%d %s: expected '%s', got '%s'", i, test.TestName, test.Result, result)
			}
		}
	})

	t.Run("GetConfig", func(t *testing.T) {
		type data struct {
			TestName string
			CfgName  string
			Err      error
			Result   Cfg
		}

		valid := Cfg{
			Name:       "valid",
			Hostname:   "hostname",
			Port:       "port",
			Abs:        "abs",
			Username:   "username",
			PrivateKey: "private_key",
		}

		tests := []data{
			{"Success", "valid", nil, valid},
			{"Noop", "noop", fmt.Errorf("open %s: no such file or directory", filepath.Join(configCtlHome, "configs", "noop", "config.json")), Cfg{}},
		}

		for i, test := range tests {
			var result Cfg
			err := GetConfig(test.CfgName, &result)
			if err != nil && test.Err == nil {
				t.Errorf("Unexpected error happend: %s\n", err)
				continue
			}

			if err == nil && test.Err != nil {
				t.Errorf("Unexpected error happend. err expected non-nil, but got nil")
				continue
			}

			if test.Err != nil {
				if err.Error() != test.Err.Error() {
					t.Errorf("Test #%d %s (Error): expected '%s', got '%s'", i, test.TestName, test.Err, err)
				}
				continue
			}

			if !reflect.DeepEqual(result, test.Result) {
				t.Errorf("Test #%d %s: expected '%+v', got '%+v'", i, test.TestName, test.Result, result)
			}
		}

	})
}
