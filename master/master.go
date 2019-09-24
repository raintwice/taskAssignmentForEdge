package main

import (
 //   "log"
//   "net"

//    "golang.org/x/net/context"
//	"google.golang.org/grpc"
	"taskAssignmentForEdge/master/mastergrpc"
  "taskAssignmentForEdge/master/taskmgt"
//    pb "taskAssignmentForEdge/proto"
    "sync"
)


func main() {
    var wg sync.WaitGroup
    wg.Add(1)
    go mastergrpc.BootupGrpcServer()
    taskmgt.ReadTaskList("tasklist.json")
    wg.Wait()
}
