// Copyright © 2018 himetani
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/himetani/configctl/workspace"
	"github.com/spf13/cobra"
)

var (
	hostname   string
	port       string
	abs        string
	username   string
	privateKey string
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [jobName]",
	Short: "Add new job",
	Long:  `Add new job`,
}

func init() {
	addCmd.RunE = add
	addCmd.Flags().StringVar(&hostname, "hostname", "", "hostname")
	addCmd.Flags().StringVar(&port, "port", "2222", "port number")
	addCmd.Flags().StringVar(&abs, "abs", "", "absolutely path")
	addCmd.Flags().StringVar(&username, "username", "", "username")
	addCmd.Flags().StringVar(&privateKey, "private-key", filepath.Clean(os.Getenv("CONFIGCTL_TEST_PRIVATE_KEY")), "ssh private key path")
	RootCmd.AddCommand(addCmd)
}

func add(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Arguments are invalid")
	}

	silent(cmd)

	name := args[0]

	job := &workspace.Job{
		Name:        name,
		Hostname:    hostname,
		Port:        port,
		Abs:         abs,
		Username:    username,
		PrivateKey:  privateKey,
		LastUpdated: time.Now(),
		LatestIdx:   0,
	}

	return workspace.CreateJob(job)
}
