package main

import (
	"log"
	"os"
	"strconv"

	//"os"
	//"strconv"
	"sync"
	"taskAssignmentForEdge/node/connect"
)

const (
    Maddress  = "localhost" //master ip
    Mport  = 50051		    //master port
	Sport = 50052
)

//Usage:L ./node master-ip node-port
func main(){
	var madd string = Maddress
	var sport int = Sport

	var _ error
	switch len(os.Args) {
	case 2:
		madd = os.Args[1]
		log.Printf("args 1 : %s", madd)
	case 3:
		madd = os.Args[1]
		log.Printf("args 1 : %s", madd)
		sport, _ = strconv.Atoi(os.Args[2])
		log.Printf("args 2 : %s, %d", os.Args[2], sport)
	default:
	}

	node := connect.NewNode(madd, Mport, sport)
	node.Init()
	var wg sync.WaitGroup
    wg.Add(1)
    go node.Join()
	wg.Add(1)
	go node.StartHeartbeatSender()
    wg.Wait()
}
