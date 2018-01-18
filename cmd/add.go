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

	"github.com/himetani/configctl/cfg"
	"github.com/himetani/configctl/workspace"
	"github.com/spf13/cobra"
)

var (
	hostname string
	abs      string
	username string
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add new configuration operation",
	Long:  `Add new configuration operation`,
}

func init() {
	addCmd.RunE = add
	addCmd.Flags().StringVar(&hostname, "hostname", "", "hostname")
	addCmd.Flags().StringVar(&abs, "abs", "", "absolutely path")
	addCmd.Flags().StringVar(&username, "username", "", "username")
	RootCmd.AddCommand(addCmd)
}

func add(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Arguments are invalid")
	}

	silent(cmd)

	name := args[0]

	cfg := &cfg.Cfg{
		Name:     name,
		Hostname: hostname,
		Abs:      abs,
		Username: username,
	}

	/*
		if err := workspace.CreateConfig(name); err != nil {
			return err
		}
	*/

	return workspace.CreateConfig(cfg)
}
