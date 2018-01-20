package workspace

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestInit(t *testing.T) {
	setup()
	defer teardown()

	configCtlHome := os.Getenv("CONFIGCTL_HOME")

	if err := Init(); err != nil {
		t.Errorf("Unexpected error happend: %s\n", err)
	}

	if _, err := os.Stat(configCtlHome); os.IsNotExist(err) {
		t.Errorf("$CONFIGCTL_HOME dir is not created. %s\n", err)
	}

	if _, err := os.Stat(filepath.Join(configCtlHome, "configs")); os.IsNotExist(err) {
		t.Errorf("$CONFIGCTL_HOME/configs is not created. %s\n", err)
	}
}

func setup() {
	testHome := os.Getenv("TEST_CONFIGCTL_HOME")
	os.Setenv("CONFIGCTL_HOME", testHome)
}

func teardown() {
	if err := os.RemoveAll(os.Getenv("CONFIGCTL_HOME")); err != nil {
		log.Printf("Unexpected error happened at teardown: %s\n", err)
	}
}
