package workspace

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

var configCtlHome string

// Init is function initialize workspace
func Init(path string) error {
	if path == "" {
		home, err := homedir.Dir()
		if err != nil {
			return err
		}
		configCtlHome = filepath.Join(home, ".configctl")
	} else {
		configCtlHome = path
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

// CreateConfig create new config
func CreateConfig(cfg *Cfg) error {
	configs := getConfigs()

	for _, c := range configs {
		if c == cfg.Name {
			return fmt.Errorf("[ERROR] %s is already created", cfg.Name)
		}
	}

	cfgPath := filepath.Join(configCtlHome, "configs", cfg.Name)
	if err := os.Mkdir(cfgPath, 0777); err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(cfgPath, "config.json"))
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	if err = encoder.Encode(cfg); err != nil {
		return nil
	}

	fmt.Printf("[INFO] %s is added. Confituration path is %s\n", cfg.Name, cfgPath)
	return nil
}

func getConfigs() (configs []string) {
	configPaths, _ := ioutil.ReadDir(filepath.Join(configCtlHome, "configs"))
	for _, c := range configPaths {
		configs = append(configs, filepath.Base(c.Name()))
	}

	return configs
}

// GetConfig returns configuration of operation
func GetConfig(name string, out interface{}) error {
	file, err := os.Open(filepath.Join(configCtlHome, "configs", name, "config.json"))
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(file)

	return decoder.Decode(out)
}
