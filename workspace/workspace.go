package workspace

import (
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

// Init is function initialize workspace
func Init(configCtlHome string) error {
	if configCtlHome == "" {
		home, err := homedir.Dir()
		if err != nil {
			return err
		}
		configCtlHome = filepath.Join(home, ".configctl")
	}

	if _, err := os.Stat(configCtlHome); os.IsNotExist(err) {
		if err := os.Mkdir(configCtlHome, 0777); err != nil {
			return err
		}
	}

	configs := filepath.Join(configCtlHome, "configs")
	if _, err := os.Stat(configs); os.IsNotExist(err) {
		if err := os.Mkdir(configs, 0777); err != nil {
			return err
		}
	}

	return nil
}
