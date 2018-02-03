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
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/himetani/configctl/workspace"
	"github.com/spf13/cobra"
)

var (
	hostStr    string
	port       string
	abs        string
	username   string
	privateKey string
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [jobName]",
	Short: "Create new job",
	Long:  `Create new job`,
}

func init() {
	createCmd.RunE = create
	createCmd.Flags().StringVar(&hostStr, "hosts", "", "hostnames. Multiple hosts are separeted by colon. (ex) host1.co.jp:host2.co.jp")
	createCmd.Flags().StringVar(&port, "port", "", "port number")
	createCmd.Flags().StringVar(&abs, "abs", "", "absolutely path. The target configuration file absolutely path.")
	createCmd.Flags().StringVar(&username, "username", "", "username")
	createCmd.Flags().StringVar(&privateKey, "private-key", "", "ssh private key path")
	RootCmd.AddCommand(createCmd)
}

func create(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Arguments are invalid")
	}

	silent(cmd)

	var hosts []string
	name := args[0]

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Job name: %s \n", name)

	if hostStr == "" {
		for {
			fmt.Println("Multiple hosts are separeted by colon. (ex) host1.co.jp:host2.co.jp")
			fmt.Print("Hosts: ")
			hostStr, _ = reader.ReadString('\n')
			hostStr = strings.Replace(hostStr, "\n", "", -1)
			if hostStr == "" {
				continue
			}

			hosts = strings.Split(hostStr, ":")
			fmt.Printf("Hosts = %v \n", hosts)
			break
		}
	}

	if port == "" {
		for {
			fmt.Print("Port (default=22): ")
			port, _ = reader.ReadString('\n')
			port = strings.Replace(port, "\n", "", -1)
			if port == "" {
				port = "22"
			}
			portInt, err := strconv.Atoi(port)
			if err != nil {
				fmt.Printf("Invalid Port number. Port = %s \n", port)
				continue
			}

			if (1024 < portInt && portInt < 65535) || portInt == 22 {
				fmt.Printf("Port = %s \n", port)
				break
			} else {
				fmt.Printf("Invalid Port number. Port = %s \n", port)
			}
		}
	}

	if abs == "" {
		for {
			fmt.Print("Config file Abs path: ")
			abs, _ = reader.ReadString('\n')
			abs = strings.Replace(abs, "\n", "", -1)
			if abs == "" {
				continue
			}

			fmt.Printf("Abs path = %s \n", abs)
			break
		}
	}

	if username == "" {
		for {
			fmt.Print("username : ")
			username, _ = reader.ReadString('\n')
			username = strings.Replace(username, "\n", "", -1)
			if username == "" {
				continue
			}
			fmt.Printf("username = %s \n", username)
			break
		}
	}

	if privateKey == "" {
		for {
			fmt.Print("privateKey : ")
			privateKey, _ = reader.ReadString('\n')
			privateKey = strings.Replace(privateKey, "\n", "", -1)

			if _, err := os.Stat(privateKey); os.IsNotExist(err) {
				fmt.Printf("Invalid privateKey path. %s\n", err)
				continue
			} else {
				fmt.Printf("privateKey = %s \n", privateKey)
				break
			}
		}
	}

	fmt.Printf("  Job Name     \t: %s\n", name)
	fmt.Printf("  Hostname     \t: %s\n", hosts)
	fmt.Printf("  Port         \t: %s\n", port)
	fmt.Printf("  AbsolutePath \t: %s\n", abs)
	fmt.Printf("  Username     \t: %s\n", username)
	fmt.Printf("  PrivateKey   \t: %s\n", privateKey)

	for {
		fmt.Print("Do you create new job? (y/n): ")
		ans, _ := reader.ReadString('\n')
		ans = strings.Replace(ans, "\n", "", -1)
		if ans == "y" {
			job := &workspace.Job{
				Name:        name,
				Hosts:       hosts,
				Port:        port,
				Abs:         abs,
				Username:    username,
				PrivateKey:  privateKey,
				LastUpdated: time.Now(),
				LatestIdx:   0,
			}
			return workspace.CreateJob(job)
		} else if ans == "n" {
			fmt.Println("Job creation aborted.")
			return nil
		} else {
			continue
		}
	}
}
