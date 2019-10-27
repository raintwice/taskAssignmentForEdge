package main

import (
	"flag"
	"fmt"
	//"log"
	"os"
	//"strconv"

	//"os"
	//"strconv"
	"sync"
	"taskAssignmentForEdge/node/connect"
	//"taskAssignmentForEdge/common"
)

var (
	h bool
	masterip string
	masterport int
	nodeport int
)

func init() {
	flag.BoolVar(&h, "h", false, "print help information")
	flag.StringVar(&masterip, "mip", "127.0.0.1", "set the listener ip of master")
	flag.IntVar(&masterport, "mport", 50051, "set the listener port of master")
	flag.IntVar(&nodeport, "nport", 50052, "set the listener port of node")
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: node [-h] [-mip masterip] [-mport masterport] [-nport nodeport]
Options:
`)
	flag.PrintDefaults()
}

func main(){
	flag.Parse()
	if h {
		usage()
		//flag.Usage()
		return
	}

	node := connect.NewNode(masterip, masterport, nodeport)
	node.InitConnection()
	node.Join()
	var wg sync.WaitGroup
    wg.Add(1)
	go node.StartHeartbeatSender(&wg)
	wg.Add(1)
	go node.StartRecvServer(&wg)
	wg.Add(1)
	go node.StartPool(3, &wg)
    wg.Wait()
}
