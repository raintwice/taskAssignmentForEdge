package main

import (
	"flag"
	"fmt"
	"taskAssignmentForEdge/common"

	//"log"
	"os"
	//"strconv"

	//"os"
	//"strconv"
	"sync"
	"taskAssignmentForEdge/node/connect"
)

var (
	h bool
	masterIP string
	masterPort int
	nodePort int
	poolCap  int
	bandWidth float64
	machineType int
	groupIndex int
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
	flag.IntVar(&poolCap, "poolcap", 3, "set the capacity of the task pool")
	flag.Float64Var(&bandWidth, "bw", 5.0, "set the bandwidth of the link between master and this node(Mbps)")
	flag.IntVar(&machineType, "machineType", common.MachineType_Simualted, "set the machine type of the edge node ")
	flag.IntVar(&startTime, "startTime", 1, "set the duration to join after booting up")
	flag.IntVar(&groupIndex, "groupIndex", common.GroupIndex_Simulated, "set the group index of the edge node")
	flag.IntVar(&pscTimeAvg, "pscTimeAvg", 15, "set the average of presence time of the edge node(in minutes) ")
	flag.IntVar(&pscTimeSig, "pscTimeSig", 15*60*0.1, "set the standard deviation of presence time of the edge node(in seconds) ")
	flag.Float64Var(&avl, "avl", 0.80, "set the availability of the edge node(in minutes) ")
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: node [-h] [-mip masterip] [-mport masterport] [-nport nodeport] [-poolcap cap_of_pool] [-bw bandwidth] 
[-g type_of_machine] [-groupIndex index_of_group] [-pscTimeAvg presence_time_avg] [-pscTimeSig presence_time_sig] [-avl avl]

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
	node.SetNodePara(bandWidth, machineType, startTime, poolCap)
	node.SetNetworkPara(groupIndex, pscTimeAvg, pscTimeSig, avl)
	node.InitConnection()
	//node.Join()
	var wg sync.WaitGroup
	if pscTimeSig == 0 {
		node.Join()
	} else {
		wg.Add(1)
		go node.StartNetworkManager(&wg)
	}
    wg.Add(1)
	go node.StartHeartbeatSender(&wg)
	wg.Add(1)
	go node.StartRecvServer(&wg)
    wg.Wait()
}
