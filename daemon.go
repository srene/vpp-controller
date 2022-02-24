
package main

import (
	"context"
	"git.fd.io/govpp.git/api"
	interfaces "git.fd.io/govpp.git/binapi/interface"
	"git.fd.io/govpp.git/binapi/vpe"
	"git.fd.io/govpp.git/binapi/ip"
	"git.fd.io/govpp.git/binapi/ip_types"
	"git.fd.io/govpp.git/binapi/sr"

	"regexp"
	"strconv"
	"strings"
	"time"

	//"flag"
	"fmt"
	"git.fd.io/govpp.git/core"
	"log"
	"os"
	"os/signal"

	"github.com/spf13/cobra"

	"git.fd.io/govpp.git"

)
var errs []error

//var log = logging.Logger("cmd")

/*var (
	sockAddr = flag.String("sock", socketclient.DefaultSocketName, "Path to VPP binary API socket file")
)*/

// getCmd represents the get command
func InitDaemonCmd(comm *Common) *cobra.Command {
	return &cobra.Command{
		Use:   "daemon",
		Short: "Run controller daemon",
		Long: ``,
		Run: func(cmd *cobra.Command, args []string) {
			datahopCli := `
 		 _     _____ _____ _   _    _____            _             _ _
		| |   |_   _/ ____| \ | |  / ____|          | |           | | |          
		| |__   | || |    |  \| | | |     ___  _ __ | |_ _ __ ___ | | | ___ _ __ 
		| '_ \  | || |    | .   | | |    / _ \| '_ \| __| '__/ _ \| | |/ _ \ '__|
		| | | |_| || |____| |\  | | |___| (_) | | | | |_| | | (_) | | |  __/ |
		|_| |_|_____\_____|_| \_|  \_____\___/|_| |_|\__|_|  \___/|_|_|\___|_|
	`
			conn, connEv, err := govpp.AsyncConnect("/run/vpp/api2.sock", core.DefaultMaxReconnectAttempts, core.DefaultReconnectInterval)
			if err != nil {
				log.Fatalln("ERROR:", err)
			}
			defer conn.Disconnect()

			// wait for Connected event
			select {
			case e := <-connEv:
				if e.State != core.Connected {
					log.Fatalln("ERROR: connecting to VPP failed:", e.Error)
				}
			}


			// check compatibility of used messages
			ch, err := conn.NewAPIChannel()

			//comm.Connection = conn
			comm.Channel = ch

			comm.Routes = make(map[string]string)
			fmt.Println(datahopCli)
			fmt.Println("VPP controller daemon running")

			if err != nil {
				log.Fatalln("ERROR: creating channel failed:", err)
			}
			defer ch.Close()
			if err := ch.CheckCompatiblity(vpe.AllMessages()...); err != nil {
				log.Fatal(err)
			}
			if err := ch.CheckCompatiblity(interfaces.AllMessages()...); err != nil {
				log.Fatal(err)
			}

			// process errors encountered during the example
			defer func() {
				if len(errs) > 0 {
					fmt.Printf("finished with %d errors\n", len(errs))
					os.Exit(1)
				} else {
					fmt.Println("finished successfully")
				}
			}()

			// send and receive messages using stream (low-low level API)
			stream, err := conn.NewStream(context.Background(),
				core.WithRequestSize(50),
				core.WithReplySize(50),
				core.WithReplyTimeout(2*time.Second))
			if err != nil {
				panic(err)
			}
			defer func() {
				if err := stream.Close(); err != nil {
					log.Fatalln(err, "closing the stream")
				}
			}()
			comm.Stream = stream

			result := make(chan int)
			//fmt.Println("1")
			//t := time.NewTicker(3 * time.Second)

			var pkt int = 0
			createPolicy(ch)
			for {
				go ipRouteDumpStream(stream,result)
				//fmt.Println("2")
				pktTemp := <- result
				if pkt != pktTemp {
					if _, ok := comm.Routes["b001::"]; ok {
						fmt.Println("Creating policy for prefix b001:: to 2::2")
						createPolicy(ch)
					} else {
						fmt.Println("No route for b001:: prefix")
					}
				}
				pkt = pktTemp
				time.Sleep(100 * time.Millisecond)
			}

			var sigChan chan os.Signal
			sigChan = make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt)
			for {
				select {
				case <-sigChan:
					//fmt.Println("cancel")
					comm.Cancel()
					return
				case <-comm.Context.Done():
					return
				}
			}
		},
	}
}

func ipRouteDumpStream(stream api.Stream, out chan<- int)  {
	//fmt.Printf("Dumping IP routes\n")

	if err := stream.SendMsg(&ip.IPRouteDump{
		ip.IPTable{
			IsIP6: true,
		},
	}); err != nil {
		logError(err, "IPRouteDump sending message")

	}
	if err := stream.SendMsg(&vpe.ControlPing{}); err != nil {
		logError(err, "ControlPing sending sending message")
	}

Loop:
	for {
		msg, err := stream.RecvMsg()
		if err != nil {
			logError(err, "IPAddressDump receiving message ")
		}

		switch msg.(type) {
		case *ip.IPRouteDetails:
		//	fmt.Printf("%+v\n",msg)
			matched, _ := regexp.MatchString("::/0", fmt.Sprintf("%+v\n",msg))
			if matched {
				//fmt.Println(msg)
				route := getPacket(fmt.Sprintf("%+v\n",msg))
				out <- route
				//fmt.Println(route)
				//return route
			}

			//break Loop
		case *vpe.ControlPingReply:
			//fmt.Printf(" - ControlPingReply: %+v\n", msg)
			break Loop

		default:
			logError(err, "unexpected message")
		}
	}


	/*fmt.Println("OK")
	fmt.Println()*/
}


func logError(err error, msg string) {
	fmt.Printf("ERROR: %s: %v\n", msg, err)
	errs = append(errs, err)
}

func getPacket(s string) int {
	//_, i := utf8.DecodeRuneInString(s)
	//s[i:]
	sp := strings.Split(s, " ")
	//fmt.Println(sp[4])
	value := strings.Replace(sp[4], "Packets:", "", -1)
	result,_ :=strconv.Atoi(value)
	return result
}


func createPolicy(ch api.Channel) {
	ip_encap,_:= ip_types.ParseIP6Address("1::1")
	sr_encap := &sr.SrSetEncapSource{
		EncapsSource: ip_encap,
	}
	sr_encap_reply := &sr.SrSetEncapSourceReply{}
	err := ch.SendRequest(sr_encap).ReceiveReply(sr_encap_reply)

	if err != nil {
		fmt.Errorf("create_encap: %w\n", err)
	}

	//fmt.Printf("create_encap: ret val %d\n",
	//	int(sr_encap_reply.Retval))

	time.Sleep(2 * time.Second)
	ip,_:= ip_types.ParseIP6Address("1::1:999")
	ip_sid,_:= ip_types.ParseIP6Address("2::2")
	sids :=[16]ip_types.IP6Address{ip_sid}
	sr_create := &sr.SrPolicyAdd{BsidAddr: ip,
		IsEncap: true,
		Sids: sr.Srv6SidList{NumSids:1,Sids: sids},
	}
	/*sr_create := &sr.SrPolicyAdd{
	                        BsidAddr: ip,
	//                              Weight:   1,
	                                IsEncap:  true,
	//                              IsSpray:  false,
	                              FibTable: 1,
	//                              Sids:     sr.Srv6SidList{Weight: 1},
	        }*/
	sr_create_reply := &sr.SrPolicyAddReply{}

	err = ch.SendRequest(sr_create).ReceiveReply(sr_create_reply)

	if err != nil {
		fmt.Errorf("create_policy: %w\n", err)
	}

	//fmt.Printf("create_policy: ret val %d\n",
	//	int(sr_create_reply.Retval))

	ip_steering,_:=ip_types.ParsePrefix("b001::/64")
	sr_steering := &sr.SrSteeringAddDel{IsDel: false,
		BsidAddr: ip,
		Prefix: ip_steering,
		TrafficType: 6}
	sr_steering_reply := &sr.SrSteeringAddDelReply{}
	err = ch.SendRequest(sr_steering).ReceiveReply(sr_steering_reply)

	if err!= nil {
		fmt.Errorf("create_steering %w\n", err)
	}

	//fmt.Printf("create steering: ret val %d\n",
	//	int(sr_steering_reply.Retval))
}
