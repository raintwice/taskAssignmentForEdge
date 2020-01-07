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
	masterIP string
	masterPort int
	nodePort int
	poolCap  int
	bandWidth float64
	machineType int
	pscTimeAvg  int
	pscTimeSig int
	avl float64
	startTime int
)

func init() {
	flag.BoolVar(&h, "h", false, "print help information")
	flag.StringVar(&masterIP, "mip", "127.0.0.1", "set the listener ip of master")
	flag.IntVar(&masterPort, "mport", 50051, "set the listener port of master")
	flag.IntVar(&nodePort, "nport", 50052, "set the listener port of node")
	flag.IntVar(&poolCap, "poolcap", 5, "set the capacity of the task pool")
	flag.Float64Var(&bandWidth, "bandwidth", 5.0, "set the bandwidth of the link between master and this node ")
	flag.IntVar(&machineType, "machineType", connect.MachineType_Simualted, "set the machine type of the edge node ")
	flag.IntVar(&pscTimeAvg, "pscTimeAvg", 15, "set the average of presence time of the edge node(in minutes) ")
	flag.IntVar(&pscTimeSig, "pscTimeSig", 15, "set the standard deviation of presence time of the edge node(in seconds) ")
	flag.Float64Var(&avl, "avl", 0.99, "set the availability of the edge node(in minutes) ")
	flag.IntVar(&startTime, "startTime", 3, "set the duration to join after booting up")
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: node [-h] [-mip masterip] [-mport masterport] [-nport nodeport] [-poolcap cap_of_pool] [-bw bandwidth]
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

	node := connect.NewNode(masterIP, masterPort, nodePort)
	node.SetParameters(bandWidth, machineType, pscTimeAvg, pscTimeSig, avl, startTime)
	node.InitConnection()
	//node.Join()
	var wg sync.WaitGroup
	wg.Add(1)
	go node.StartNetworkManager(&wg)
    wg.Add(1)
	go node.StartHeartbeatSender(&wg)
	wg.Add(1)
	go node.StartRecvServer(&wg)
	wg.Add(1)
	go node.StartPool(poolCap, &wg)
    wg.Wait()
}
