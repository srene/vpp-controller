/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"bytes"
	"context"
	"fmt"
	"log"

	uds "github.com/asabya/go-ipc-uds"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	//"github.com/srene/vpp-controller/cmd"
	//"github.com/srene/vpp-controller/common"
	"path/filepath"
	"strings"

	"os"
)

const (
	argSeparator = "$^~@@*"
)

var (
	rootCmd = &cobra.Command{
		Use:   "vpp-controller",
		Short: "This is vpp-controller cli client",
		Long: `
The VPP controller CLI client gives access to the controller through a CLI Interface.
		`,
	}
	sockPath = "uds.sock"
	//log      = logging.Logger("cmd")
)

func init() {
//	logging.SetLogLevel("uds", "Debug")
//	logging.SetLogLevel("cmd", "Debug")
}



func main() {

	ctx, cancel := context.WithCancel(context.Background())
	/*home, err := os.UserHomeDir()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	root := filepath.Join(home, repo.Root)
	err = repo.Init(root, "0")
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}*/

	comm := &Common{
	//	Root:    root,
		Context: ctx,
		Cancel:  cancel,
	}

	//rootCmd.PersistentFlags().BoolP("json", "j", false, "json output")
	//rootCmd.PersistentFlags().BoolP("pretty", "p", false, "pretty json output")
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	var allCommands []*cobra.Command
	allCommands = append(
		allCommands,
		InitDaemonCmd(comm),
		InitAddCmd(comm),
		InitListCmd(comm),
		/*cmd.InitInfoCmd(comm),
		cmd.InitStopCmd(comm),
		cmd.InitAddCmd(comm),
		cmd.InitIndexCmd(comm),
		cmd.InitRemoveCmd(comm),
		cmd.InitGetCmd(comm),
		cmd.InitVersionCmd(comm),
		cmd.InitMatrixCmd(comm),
		cmd.InitializeDocCommand(comm),
		cmd.InitGetCmd(comm),
		cmd.InitCompletionCmd(comm),*/
	)

	for _, i := range allCommands {
		rootCmd.AddCommand(i)
	}
	// check help flag
	for _, v := range os.Args {
		if v == "-h" || v == "--help" {
			log.Println("Executing help command")
			rootCmd.Execute()
			return
		}
	}

	socketPath := filepath.Join("/tmp", sockPath)
	/*if !uds.IsIPCListening(socketPath) {
		r, err := repo.Open(root)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		defer r.Close()
		comm.Repo = r
	}*/
	if len(os.Args) > 1 {
		if os.Args[1] != "daemon" && uds.IsIPCListening(socketPath) {
			opts := uds.Options{
				SocketPath: filepath.Join("/tmp", sockPath),
			}
			r, w, c, err := uds.Dialer(opts)
			if err != nil {
				log.Fatalln(err)
				goto Execute
			}
			defer c()
			err = w(strings.Join(os.Args[1:], argSeparator))
			if err != nil {
				log.Fatalln(err)
				os.Exit(1)
			}
			v, err := r()
			if err != nil {
				log.Fatalln(err)
				os.Exit(1)

			}
			fmt.Println(v)
			return
		}
		if os.Args[1] == "daemon" {
			if uds.IsIPCListening(socketPath) {
				fmt.Println("VPP controller daemon is already running")
				return
			}
			_, err := os.Stat(filepath.Join("/tmp", sockPath))
			if !os.IsNotExist(err) {
				err := os.Remove(filepath.Join("/tmp", sockPath))
				if err != nil {
					log.Fatalln(err)
					os.Exit(1)
				}
			}
			opts := uds.Options{
				SocketPath: filepath.Join("/tmp", sockPath),
			}
			in, err := uds.Listener(context.Background(), opts)
			if err != nil {
				log.Fatalln(err)
				os.Exit(1)
			}
			go func() {
				for {
					client := <-in
					go func() {
						for {
							ip, err := client.Read()
							if err != nil {
								break
							}
							if len(ip) == 0 {
								break
							}
							commandStr := string(ip)
							//log.Println("run command :", commandStr)
							var (
								childCmd *cobra.Command
								flags    []string
							)
							command := strings.Split(commandStr, argSeparator)
							if rootCmd.TraverseChildren {
								childCmd, flags, err = rootCmd.Traverse(command)
							} else {
								childCmd, flags, err = rootCmd.Find(command)
							}
							if err != nil {
								err = client.Write([]byte(err.Error()))
								if err != nil {
									log.Fatalln("Write error", err)
									client.Close()
								}
								break
							}
							childCmd.Flags().VisitAll(func(f *pflag.Flag) {
								err := f.Value.Set(f.DefValue)
								if err != nil {
									log.Fatalln("Unable to set flags ", childCmd.Name(), f.Name, err.Error())
								}
							})
							if err := childCmd.Flags().Parse(flags); err != nil {
								log.Fatalln("Unable to parse flags ", err.Error())
								err = client.Write([]byte(err.Error()))
								if err != nil {
									log.Fatalln("Write error", err)
									client.Close()
								}
								break
							}
							outBuf := new(bytes.Buffer)
							childCmd.SetOut(outBuf)
							if childCmd.Args != nil {
								if err := childCmd.Args(childCmd, flags); err != nil {
									err = client.Write([]byte(err.Error()))
									if err != nil {
										log.Fatalln("Write error", err)
										client.Close()
									}
									break
								}
							}
							if childCmd.PreRunE != nil {
								if err := childCmd.PreRunE(childCmd, flags); err != nil {
									err = client.Write([]byte(err.Error()))
									if err != nil {
										log.Fatalln("Write error", err)
										client.Close()
									}
									break
								}
							} else if childCmd.PreRun != nil {
								childCmd.PreRun(childCmd, command)
							}

							if childCmd.RunE != nil {
								if err := childCmd.RunE(childCmd, flags); err != nil {
									err = client.Write([]byte(err.Error()))
									if err != nil {
										log.Fatalln("Write error", err)
										client.Close()
									}
									break
								}
							} else if childCmd.Run != nil {
								childCmd.Run(childCmd, flags)
							}

							out := outBuf.Next(outBuf.Len())
							outBuf.Reset()
							err = client.Write(out)
							if err != nil {
								log.Fatalln("Write error", err)
								client.Close()
								break
							}
						}
					}()
				}
			}()
		}
	}
Execute:
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}


