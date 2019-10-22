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
    go ms.StartGrpcServer()
    wg.Add(1)
    go ms.StartHeartbeatChecker()
    //taskmgt.ReadTaskList("tasklist.json")
    wg.Wait()
}
