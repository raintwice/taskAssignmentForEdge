package main

import (
//    "log"
    //"os"
  "sync"
  "math/rand"
  "time"
  "fmt"
  //"strconv"
//    "golang.org/x/net/context"
//	"google.golang.org/grpc"
	"taskAssignmentForEdge/node/nodegrpc"
//    pb "taskAssignmentForEdge/proto"
)


func main(){
	var wg sync.WaitGroup
    wg.Add(1)
    // select a random port in 40000 - 50000 and listen to it
    rand.Seed(time.Now().UnixNano())
    nodeport := 40000 + rand.Intn(10000)
    ipAddr := fmt.Sprintf("localhost:%d", nodeport)
    //p32 := int32(nodeport)
    go nodegrpc.BootNodeGrpc(ipAddr)
    // join in the master
    go nodegrpc.Join(ipAddr)
    wg.Wait()
}
