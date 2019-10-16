package main

import (
 //   "log"
//   "net"

//    "golang.org/x/net/context"
//	"google.golang.org/grpc"
	"taskAssignmentForEdge/master/common"
  //"taskAssignmentForEdge/master/taskmgt"
//    pb "taskAssignmentForEdge/proto"
    "sync"
)


func main() {
    ms := common.NewMaster()
    ms.Init()
    var wg sync.WaitGroup
    wg.Add(1)
    go ms.BootupGrpcServer()
    //taskmgt.ReadTaskList("tasklist.json")
    wg.Wait()
}
