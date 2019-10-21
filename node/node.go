package main

import (
//    "log"
    //"os"
    "sync"

//    "golang.org/x/net/context"
//	"google.golang.org/grpc"
	"taskAssignmentForEdge/node/common"
//    pb "taskAssignmentForEdge/proto"
)

const (
    Maddress  = "localhost" //master ip
    Mport  = 50051		    //master port
)

func main(){
	node := common.NewNode(Maddress, Mport)
	node.Init()
	var wg sync.WaitGroup
    wg.Add(1)
    go  node.Join()
    wg.Wait()
}
