package cmd

import (
	"fmt"
	"log"
	"github.com/spf13/cobra"
	"github.com/srene/vpp-controller/common"
	"github.com/srene/vpp-controller/out"
)

// getCmd represents the get command
func InitListCmd(comm *common.Common) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "list routes",
		Long: ``,
		RunE: func(cmd *cobra.Command, args []string) error{
			fmt.Println("List command")
			err := out.Print(cmd, nil,parseFormat(cmd))
			if err != nil {
				log.Fatalln("Unable to get config ", err)
				return err
			}

			return nil

		},
	}
}
