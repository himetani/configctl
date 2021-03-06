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

	"github.com/himetani/configctl/workspace"
	"github.com/spf13/cobra"
)

// applyHisoryCmd represents the applyHisory command
var applyHistoryCmd = &cobra.Command{
	Use:   "applyHistory [jobName] [idx]",
	Short: "Show history of apply execution",
	Long:  `Show history of apply execution`,
}

func init() {
	applyHistoryCmd.RunE = applyHistory
	RootCmd.AddCommand(applyHistoryCmd)

}

func applyHistory(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return errors.New("Arguments are invalid")
	}

	silent(cmd)
	name := args[0]
	idx := args[1]

	return workspace.ShowHistory(name, idx)
}
