package main

import (
//    "log"
    //"os"
    "sync"

//    "golang.org/x/net/context"
//	"google.golang.org/grpc"
	"taskAssignmentForEdge/node/nodegrpc"
//    pb "taskAssignment/proto"
)


func main(){
	var wg sync.WaitGroup
    wg.Add(1)
    go nodegrpc.Join()
    wg.Wait()
}
