package workspace

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestInit(t *testing.T) {
	if err := setup(); err != nil {
		t.Fatal(err)
	}

	configCtlHome := os.Getenv("CONFIGCTL_HOME")

	if err := Init(); err != nil {
		t.Errorf("Unexpected error happend: %s\n", err)
	}

	defer teardown()

	if _, err := os.Stat(configCtlHome); os.IsNotExist(err) {
		t.Errorf("$CONFIGCTL_HOME dir is not created. %s\n", err)
	}

	if _, err := os.Stat(filepath.Join(configCtlHome, "configs")); os.IsNotExist(err) {
		t.Errorf("$CONFIGCTL_HOME/configs is not created. %s\n", err)
	}
}

func setup() error {
	testHome := os.Getenv("TEST_CONFIGCTL_HOME")
	if testHome == "" {
		return errors.New("Set TEST_CONFIGCTL_HOME environment variable to test")
	}
	os.Setenv("CONFIGCTL_HOME", testHome)
	return nil
}

func teardown() {
	if err := os.RemoveAll(os.Getenv("CONFIGCTL_HOME")); err != nil {
		log.Printf("Unexpected error happened at teardown: %s\n", err)
	}
}
