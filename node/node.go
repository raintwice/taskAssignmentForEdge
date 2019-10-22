package main

import (
    "sync"
	"taskAssignmentForEdge/node/connect"
)

const (
    Maddress  = "localhost" //master ip
    Mport  = 50051		    //master port
)

func main(){
	node := connect.NewNode(Maddress, Mport)
	node.Init()
	var wg sync.WaitGroup
    wg.Add(1)
    go node.Join()
	wg.Add(1)
	go node.StartHeartbeatSender()
    wg.Wait()
}
