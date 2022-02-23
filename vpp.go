package main
import "C"
import (
 "fmt"
 "log"
 "strings"
 "strconv"
// "math"
 "time"
 "context"
 "regexp"
// "sort"
 "sync"
 "git.fd.io/govpp.git"
 "git.fd.io/govpp.git/core"
 "git.fd.io/govpp.git/binapi/vpe"
 //interfaces "git.fd.io/govpp.git/binapi/interface"
 "git.fd.io/govpp.git/binapi/ip"
 "git.fd.io/govpp.git/binapi/ip_types"
 "git.fd.io/govpp.git/binapi/sr"
 "git.fd.io/govpp.git/api"


)
var mtx sync.Mutex
//export Query
func Query(address string) int {
  	mtx.Lock()
  	defer mtx.Unlock()
	//stream,_,err:=connect()
	        fmt.Println("starting connection")
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

        fmt.Println("new api channel")
        // check compatibility of used messages
        ch, err := conn.NewAPIChannel()
        if err != nil {
                log.Fatalln("ERROR: creating channel failed:", err)
        }
        defer ch.Close()
        /*if err := ch.CheckCompatiblity(vpe.AllMessages()...); err != nil {
                log.Fatal(err)
        }
        if err := ch.CheckCompatiblity(interfaces.AllMessages()...); err != nil {
                log.Fatal(err)
        }*/


        fmt.Println("new stream")
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

        fmt.Println("stream send msg")

        if err := stream.SendMsg(&ip.IPRouteDump{
                ip.IPTable{
                        IsIP6: true,
                },
        }); err != nil {
                fmt.Errorf("IPRouteDump sending message #{err}\n")
       //         return nil,nil,err
        }
        if err := stream.SendMsg(&vpe.ControlPing{}); err != nil {
                fmt.Errorf("ControlPing sending sending message #{err}\n")
         //       return nil,nil,err
        }


	if err!=nil{
		fmt.Errorf(" Error connect #{err}\n")
	} else {
	   fmt.Println("connection done")
	}
	var packets int
Loop:
	for {
		msg, err := stream.RecvMsg()
		if err != nil {
			fmt.Errorf("IPAddressDump receiving message #{err}\n")
		}

		switch msg.(type) {
		case *ip.IPRouteDetails:
			//fmt.Printf("%+v\n",msg)
			matched, _ := regexp.MatchString("::/0", fmt.Sprintf("%+v\n",msg))
			if matched {
				//fmt.Println(msg)
				route := getPacket(fmt.Sprintf("%+v\n",msg))
				packets = route
				//fmt.Println(route)
				//fmt.Println(route)
				//return route
			}

			//break Loop
		case *vpe.ControlPingReply:
			//fmt.Printf(" - ControlPingReply: %+v\n", msg)
			break Loop

		default:
			fmt.Errorf("unexpected message #{err}\n")
		}
	}

	return packets
	/*fmt.Println("OK")
	fmt.Println()*/

}

//export CreatePolicy
func CreatePolicy(prefix string,destination string) {
	mtx.Lock()
        defer mtx.Unlock()
        /*_,ch,err:=connect()
        if err!=nil{
                fmt.Errorf(" Error connect #{err}\n")
        } else {
		fmt.Println("connection done")
	}*/

	fmt.Println("starting connection")
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

        fmt.Println("new api channel")
        // check compatibility of used messages
        ch, err := conn.NewAPIChannel()
        if err != nil {
                log.Fatalln("ERROR: creating channel failed:", err)
        }
        defer ch.Close()
        /*if err := ch.CheckCompatiblity(vpe.AllMessages()...); err != nil {
                log.Fatal(err)
        }
        if err := ch.CheckCompatiblity(interfaces.AllMessages()...); err != nil {
                log.Fatal(err)
        }*/


        fmt.Println("new stream")
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

        fmt.Println("stream send msg")

        if err := stream.SendMsg(&ip.IPRouteDump{
                ip.IPTable{
                        IsIP6: true,
                },
        }); err != nil {
                fmt.Errorf("IPRouteDump sending message #{err}\n")
  //              return nil,nil,err
        }
        if err := stream.SendMsg(&vpe.ControlPing{}); err != nil {
                fmt.Errorf("ControlPing sending sending message #{err}\n")
//                return nil,nil,err
        }


	ip_encap,_:= ip_types.ParseIP6Address("1::1")
	sr_encap := &sr.SrSetEncapSource{
		EncapsSource: ip_encap,
	}
	sr_encap_reply := &sr.SrSetEncapSourceReply{}
	err = ch.SendRequest(sr_encap).ReceiveReply(sr_encap_reply)

	if err != nil {
		fmt.Errorf("create_encap: %w\n", err)
	}

	fmt.Printf("create_encap: ret val %d\n",
		int(sr_encap_reply.Retval))

	time.Sleep(2 * time.Second)
	ip,_:= ip_types.ParseIP6Address("1::1:999")
	ip_sid,_:= ip_types.ParseIP6Address(destination)
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

	fmt.Printf("create_policy: ret val %d\n",
		int(sr_create_reply.Retval))

	ip_steering,_:=ip_types.ParsePrefix(prefix)
	sr_steering := &sr.SrSteeringAddDel{IsDel: false,
		BsidAddr: ip,
		Prefix: ip_steering,
		TrafficType: 6}
	sr_steering_reply := &sr.SrSteeringAddDelReply{}
	err = ch.SendRequest(sr_steering).ReceiveReply(sr_steering_reply)

	if err!= nil {
		fmt.Errorf("create_steering %w\n", err)
	}

	fmt.Printf("create steering: ret val %d\n",
		int(sr_steering_reply.Retval))
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

func connect() (api.Stream,api.Channel,error) {

	fmt.Println("starting connection")
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

	fmt.Println("new api channel")
        // check compatibility of used messages
        ch, err := conn.NewAPIChannel()
        if err != nil {
                log.Fatalln("ERROR: creating channel failed:", err)
        }
        defer ch.Close()
        /*if err := ch.CheckCompatiblity(vpe.AllMessages()...); err != nil {
                log.Fatal(err)
        }
        if err := ch.CheckCompatiblity(interfaces.AllMessages()...); err != nil {
                log.Fatal(err)
        }*/


	fmt.Println("new stream")
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

	fmt.Println("stream send msg")

        if err := stream.SendMsg(&ip.IPRouteDump{
                ip.IPTable{
                        IsIP6: true,
                },
        }); err != nil {
                fmt.Errorf("IPRouteDump sending message #{err}\n")
		return nil,nil,err
        }
        if err := stream.SendMsg(&vpe.ControlPing{}); err != nil {
                fmt.Errorf("ControlPing sending sending message #{err}\n")
		return nil,nil,err
        }

	return stream,ch,nil


}
func main() {}

