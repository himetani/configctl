package workspace

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

var configCtlHome string

// Init is function initialize workspace
func Init() error {
	configCtlHome = os.Getenv("CONFIGCTL_HOME")

	if configCtlHome == "" {
		home, err := homedir.Dir()
		if err != nil {
			return err
		}
		configCtlHome = filepath.Join(home, ".configctl")
	} else {
		configCtlHome = filepath.Clean(configCtlHome)
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
	configs := GetConfigs()

	for _, c := range configs {
		if c == cfg.Name {
			return fmt.Errorf("[ERROR] %s is already created", cfg.Name)
		}
	}

	cfgPath := filepath.Join(configCtlHome, "configs", cfg.Name)
	if err := os.Mkdir(cfgPath, 0777); err != nil {
		return err
	}

	if err := os.Mkdir(filepath.Join(cfgPath, "history"), 0777); err != nil {
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

// GetConfigs returns the slice of config
func GetConfigs() (configs []string) {
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

// CreateTmp creates tmp dir and returns error
func CreateTmp(name string) error {
	tmpPath := filepath.Join(configCtlHome, "configs", name, "tmp")
	return os.Mkdir(tmpPath, 0777)
}

// DeleteTmp delete tmp dir and returns error
func DeleteTmp(name string) error {
	tmpPath := filepath.Join(configCtlHome, "configs", name, "tmp")
	return os.RemoveAll(tmpPath)
}

// PutTmp creates file in tmp dir and returns error
func PutTmp(config, name string, reader io.Reader) error {
	tmpPath := filepath.Join(configCtlHome, "configs", config, "tmp")
	file, err := os.Create(filepath.Join(tmpPath, name))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)

	return err
}

// TmpDiff execute vimdiff of files in tmp dir
func TmpDiff(name, before, after string) error {
	tmpPath := filepath.Join(configCtlHome, "configs", name, "tmp")
	cmd := exec.Command("vimdiff", filepath.Join(tmpPath, before), filepath.Join(tmpPath, after))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
