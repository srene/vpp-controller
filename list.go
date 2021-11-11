package main

import (
	"fmt"
	"github.com/spf13/cobra"

	"log"
)

// getCmd represents the get command
func InitListCmd(comm *Common) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "list routes",
		Long: ``,
		RunE: func(cmd *cobra.Command, args []string) error{
			fmt.Println("List command")
			err := Print(cmd, nil, parseFormat(cmd))
			if err != nil {
				log.Fatalln("Unable to get config ", err)
				return err
			}

			return nil

		},
	}
}
