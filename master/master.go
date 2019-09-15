package main

import (
 //   "log"
//   "net"

//    "golang.org/x/net/context"
//	"google.golang.org/grpc"
	"taskAssignment/master/mastergrpc"
  "taskAssignment/master/taskmgt"
//    pb "taskAssignment/proto"
    "sync"
)


func main() {
    var wg sync.WaitGroup
    wg.Add(1)
    go mastergrpc.BootupGrpcServer()
    taskmgt.ReadTaskList("tasklist.json")
    wg.Wait()
}
