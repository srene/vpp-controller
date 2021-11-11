package main

import (
	"fmt"
	"log"
	"errors"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
func InitAddCmd(comm *Common) *cobra.Command {
	return &cobra.Command{
		Use:   "add",
		Short: "add route",
		Long: ``,
		RunE: func(cmd *cobra.Command, args []string) error{
			fmt.Println("Add command")
			prefix := args[0]
			route := args[1]

			if comm.Stream == nil {
				return errors.New("daemon not running")
			}
			if _, ok := comm.Routes[prefix]; ok {
				//do something here
				log.Println("Updating route")
			} else {
				log.Println("Adding route")
			}
			comm.Routes[prefix] = route
			err := Print(cmd, nil, parseFormat(cmd))
			if err != nil {
				log.Fatalln("Unable to get config ", err)
				return err
			}

			return nil

		},
	}
}

func parseFormat(cmd *cobra.Command) Format {
	pFlag, _ := cmd.Flags().GetBool("pretty")
	jFlag, _ := cmd.Flags().GetBool("json")
	log.Println(pFlag, jFlag)
	var f Format
	if jFlag {
		f = Json
	}
	if pFlag {
		f = PrettyJson
	}
	if !pFlag && !jFlag {
		f = NoStyle
	}
	return f
}