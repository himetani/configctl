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

	"github.com/himetani/configctl/workspace"
	"github.com/spf13/cobra"
)

// describeCmd represents the describe command
var describeCmd = &cobra.Command{
	Use:   "describe [jobName]",
	Short: "describe the job configuration information",
	Long:  `describe the job configuration information`,
}

func init() {
	describeCmd.RunE = describe
	RootCmd.AddCommand(describeCmd)
}

func describe(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Arguments are invalid")
	}

	silent(cmd)
	name := args[0]

	var cfg workspace.Cfg
	if err := workspace.GetConfig(name, &cfg); err != nil {
		return err
	}

	fmt.Printf("  Name         \t: %s\n", cfg.Name)
	fmt.Printf("  Hostname     \t: %s\n", cfg.Hostname)
	fmt.Printf("  Port         \t: %s\n", cfg.Port)
	fmt.Printf("  AbsolutePath \t: %s\n", cfg.Abs)
	fmt.Printf("  Username     \t: %s\n", cfg.Username)
	fmt.Printf("  PrivateKey   \t: %s\n", cfg.PrivateKey)
	fmt.Printf("  LastUpdated  \t: %s\n", cfg.LastUpdated)
	fmt.Printf("  LatestIdx    \t: %d\n", cfg.LatestIdx)

	return nil
}
