package main

import (
    "flag"
    "fmt"
    "os"
    "sync"
   	"taskAssignmentForEdge/master/connect"
    "taskAssignmentForEdge/master/dispatch"
)

var (
    h bool
    dptindex int
)

func init() {
    flag.BoolVar(&h, "h", false, "print help information")
    flag.IntVar(&dptindex, "dispatcher", dispatch.Dispatcher_RR, "set the dispatcher")
}

func usage() {
    fmt.Fprintf(os.Stderr, `Usage: node [-h] [-dipatcher index_of_dispatcher] 
Options:
`)
    flag.PrintDefaults()
}

func main() {
    flag.Parse()
    if h {
        usage()
        //flag.Usage()
        return
    }

    ms := connect.NewMaster()
    ms.Init(dptindex)
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
