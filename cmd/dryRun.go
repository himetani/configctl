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
	"bytes"
	"errors"
	"os"

	"github.com/himetani/configctl/client"
	"github.com/himetani/configctl/workspace"
	"github.com/spf13/cobra"
)

// dryRunCmd represents the dryRun command
var dryRunCmd = &cobra.Command{
	Use:   "dryRun [jobName] [file]",
	Short: "Dry run shows the diff between current config file at server & applied config file",
	Long:  `Dry run shows the diff between current config file at server & applied config file`,
}

func init() {
	dryRunCmd.RunE = dryRun
	RootCmd.AddCommand(dryRunCmd)
}

func dryRun(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return errors.New("Arguments are invalid")
	}

	silent(cmd)
	name := args[0]

	var job workspace.Job
	if err := workspace.GetJob(name, &job); err != nil {
		return err
	}

	session, err := client.NewSession(job.Hosts[0], job.Port, job.Username, job.PrivateKey)
	if err != nil {
		return err
	}
	defer session.Close()

	data, err := session.Get(job.Abs)
	if err != nil {
		return err
	}

	bytes.NewBuffer(data)
	before := bytes.NewBuffer(data)

	after, err := os.Open(args[1])
	if err != nil {
		return err
	}
	defer after.Close()

	if err := workspace.CreateTmp(name); err != nil {
		return err
	}

	defer workspace.DeleteTmp(name)

	if err := workspace.PutTmp(name, "before", before); err != nil {
		return err
	}

	if err := workspace.PutTmp(name, "after", after); err != nil {
		return err
	}

	return workspace.TmpDiff(name, "before", "after")
}
