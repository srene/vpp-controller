package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/srene/vpp-controller/common"
	"github.com/srene/vpp-controller/out"
)

// getCmd represents the get command
func InitAddCmd(comm *common.Common) *cobra.Command {
	return &cobra.Command{
		Use:   "add",
		Short: "add route",
		Long: ``,
		RunE: func(cmd *cobra.Command, args []string) error{
			fmt.Println("Add command")
			err := out.Print(cmd, nil,parseFormat(cmd))
			if err != nil {
				log.Fatalln("Unable to get config ", err)
				return err
			}

			return nil

		},
	}
}

func parseFormat(cmd *cobra.Command) out.Format {
	pFlag, _ := cmd.Flags().GetBool("pretty")
	jFlag, _ := cmd.Flags().GetBool("json")
	log.Println(pFlag, jFlag)
	var f out.Format
	if jFlag {
		f = out.Json
	}
	if pFlag {
		f = out.PrettyJson
	}
	if !pFlag && !jFlag {
		f = out.NoStyle
	}
	return f
}