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
	"fmt"

	"github.com/himetani/configctl/client"
	"github.com/himetani/configctl/workspace"
	"github.com/spf13/cobra"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect [jobName]",
	Short: "connect to target job configuration server & confirm connectivity",
	Long:  `connect to target job configuration server & confirm connectivity`,
}

func init() {
	connectCmd.RunE = connect
	RootCmd.AddCommand(connectCmd)
}

func connect(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Arguments are invalid")
	}

	silent(cmd)
	name := args[0]

	var job workspace.Job
	if err := workspace.GetJob(name, &job); err != nil {
		return err
	}

	session, err := client.NewSession(job.Hostname, job.Port, job.Username, job.PrivateKey)
	if err != nil {
		return err
	}
	defer session.Close()

	bytes, err := session.Get(job.Abs)
	if err != nil {
		return err
	}

	fmt.Printf(string(bytes))

	return nil
}
