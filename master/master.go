package main

import (
    "sync"
	"taskAssignmentForEdge/master/connect"
)


func main() {
    ms := connect.NewMaster()
    ms.Init()
    var wg sync.WaitGroup
    wg.Add(1)
    go ms.StartGrpcServer(&wg)
    wg.Add(1)
    go ms.StartHeartbeatChecker(&wg)
    //taskmgt.ReadTaskList("tasklist.json")
    wg.Add(1)
    go ms.StartDispatcher(&wg)
    wg.Add(1)
    go ms.StartServerForClient(&wg)
    wg.Wait()
}
