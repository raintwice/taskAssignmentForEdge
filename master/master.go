package main

import (
 //   "log"
//   "net"

//    "golang.org/x/net/context"
//	"google.golang.org/grpc"
	"taskAssignment/master/mastergrpc"
//    pb "taskAssignment/proto"
    "sync"
)


func main() {
    var wg sync.WaitGroup
    wg.Add(1)
    go mastergrpc.BootupGrpcServer()
    wg.Wait()
}
