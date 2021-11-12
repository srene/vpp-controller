package main

import (
	"fmt"
	"git.fd.io/govpp.git/api"
	"git.fd.io/govpp.git/binapi/ip_types"
	"git.fd.io/govpp.git/binapi/sr"
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
			if val, ok := comm.Routes[prefix]; ok {
				//do something here
				log.Println("Updating route")
				updateRoute(comm.Channel, val)
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



func updateRoute(ch api.Channel, route string) error {
	ip,_:= ip_types.ParseIP6Address("1::1:999")
	ip_sid,_:= ip_types.ParseIP6Address(route)
	sids :=[16]ip_types.IP6Address{ip_sid}

	sr_delete := &sr.SrPolicyDel{BsidAddr: ip,
	}
	sr_delete_reply := &sr.SrPolicyDelReply{}

	err := ch.SendRequest(sr_delete).ReceiveReply(sr_delete_reply)
	if err != nil {
		log.Fatalln("ERROR: deleting policy:", err)
		return err
	}
	sr_create := &sr.SrPolicyAdd{BsidAddr: ip,
		IsEncap: true,
		Sids: sr.Srv6SidList{NumSids:1,Sids: sids},
	}
	sr_create_reply := &sr.SrPolicyAddReply{}

	err = ch.SendRequest(sr_create).ReceiveReply(sr_create_reply)

	return err

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