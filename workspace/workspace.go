package workspace

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

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

// CreateJob create new job
func CreateJob(cfg *Job) error {
	configs := GetJobs()

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
	encoder.SetIndent("", "\t")
	if err = encoder.Encode(cfg); err != nil {
		return nil
	}

	fmt.Printf("[INFO] %s is added. Confituration path is %s\n", cfg.Name, cfgPath)
	return nil
}

// UpdateJob is func to update job.json
func UpdateJob(job *Job) error {
	jobPath := filepath.Join(configCtlHome, "jobs", job.Name)
	file, err := os.Create(filepath.Join(jobPath, "job.json"))
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	return encoder.Encode(job)
}

// GetJobs returns the slice of Job
func GetJobs() (jobs []string) {
	jobPaths, _ := ioutil.ReadDir(filepath.Join(configCtlHome, "jobs"))
	for _, c := range jobPaths {
		jobs = append(jobs, filepath.Base(c.Name()))
	}

	return jobs
}

// GetJob returns job configuration of operation
func GetJob(name string, out interface{}) error {
	file, err := os.Open(filepath.Join(configCtlHome, "jobs", name, "job.json"))
	if err != nil {
		return fmt.Errorf("Can't find \"%s\" job. error: %s", name, err)
	}

	decoder := json.NewDecoder(file)

	return decoder.Decode(out)
}

// EditJob opens vim to edit job.json of target job and returns error
func EditJob(name string) error {
	jobPath := filepath.Join(configCtlHome, "jobs", name, "job.json")
	cmd := exec.Command("vim", jobPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// CreateTmp creates tmp dir and returns error
func CreateTmp(name string) error {
	tmpPath := filepath.Join(configCtlHome, "jobs", name, "tmp")
	return os.Mkdir(tmpPath, 0777)
}

// DeleteTmp delete tmp dir and returns error
func DeleteTmp(name string) error {
	tmpPath := filepath.Join(configCtlHome, "jobs", name, "tmp")
	return os.RemoveAll(tmpPath)
}

// PutTmp creates file in tmp dir and returns error
func PutTmp(config, name string, reader io.Reader) error {
	tmpPath := filepath.Join(configCtlHome, "jobs", config, "tmp")
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

// CreateHistory is func to create hisotry
func CreateHistory(name string, idx int, before, after io.Reader) error {
	histPath := filepath.Join(configCtlHome, "configs", name, "history", strconv.Itoa(idx))
	if err := os.Mkdir(histPath, 0777); err != nil {
		return err
	}

	beforef, err := os.Create(filepath.Join(histPath, "before"))
	if err != nil {
		return err
	}
	defer beforef.Close()

	afterf, err := os.Create(filepath.Join(histPath, "after"))
	if err != nil {
		return err
	}
	defer afterf.Close()

	_, err = io.Copy(beforef, before)
	if err != nil {
		return err
	}
	_, err = io.Copy(afterf, after)
	if err != nil {
		return err
	}

	return nil
}

// ShowHistory is func to open job execution history
func ShowHistory(name, idx string) error {
	jobPath := filepath.Join(configCtlHome, "jobs", name)
	if _, err := os.Stat(jobPath); err != nil {
		return fmt.Errorf("Can't find \"%s\" job. error: %s", name, err)
	}

	histPath := filepath.Join(jobPath, "history", idx)
	if _, err := os.Stat(histPath); err != nil {
		return fmt.Errorf("Can't find \"%s\" job history. idx: %s, error: %s", name, idx, err)
	}

	cmd := exec.Command("vimdiff", filepath.Join(histPath, "before"), filepath.Join(histPath, "after"))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
