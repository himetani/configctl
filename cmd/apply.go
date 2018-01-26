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
	"io/ioutil"
	"os"

	"github.com/himetani/configctl/client"
	"github.com/himetani/configctl/workspace"
	"github.com/spf13/cobra"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply configurations",
	Long:  `Apply configurations`,
}

func init() {
	applyCmd.RunE = apply
	RootCmd.AddCommand(applyCmd)
}

func apply(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return errors.New("Arguments are invalid")
	}

	silent(cmd)
	name := args[0]

	var cfg workspace.Cfg
	if err := workspace.GetConfig(name, &cfg); err != nil {
		return err
	}

	file, err := os.Open(args[1])
	if err != nil {
		return err
	}
	defer file.Close()

	applied, _ := ioutil.ReadAll(file)
	content := string(applied)

	if err := history(cfg, name, applied); err != nil {
		return err
	}

	session, err := client.NewSession(cfg.Hostname, cfg.Port, cfg.Username, cfg.PrivateKey)
	if err != nil {
		return err
	}
	defer session.Close()

	if err := session.Scp(content, cfg.Abs); err != nil {
		return err
	}

	cfg.LatestIdx++
	return workspace.UpdateConfig(&cfg)

}

func history(cfg workspace.Cfg, name string, applied []byte) error {
	session, err := client.NewSession(cfg.Hostname, cfg.Port, cfg.Username, cfg.PrivateKey)
	if err != nil {
		return err
	}
	defer session.Close()

	// Get Current File
	current, err := session.Get(cfg.Abs)
	if err != nil {
		return err
	}

	before := bytes.NewBuffer(current)
	after := bytes.NewBuffer(applied)

	return workspace.CreateHistory(name, cfg.LatestIdx, before, after)
}
