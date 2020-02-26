package main

import (
    "flag"
    "fmt"
    "os"
    "sync"
    "taskAssignmentForEdge/common"
    "taskAssignmentForEdge/master/connect"
    "taskAssignmentForEdge/master/dispatch"
)

var (
    h bool
    dptIndex int
    interval int
)

func init() {
    flag.BoolVar(&h, "h", false, "print help information")
    flag.IntVar(&dptIndex, "dispatcher", dispatch.Dispatcher_RR, "set the dispatcher")
    flag.IntVar(&interval, "interval", common.AssignInterval, "set the interval of dispatcher in ms")
}

func usage() {
    fmt.Fprintf(os.Stderr, `Usage: node [-h] [-dispatcher index_of_dispatcher] [-interval interval_of_dispatcher]
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
    ms.Init(dptIndex, interval)
    var wg sync.WaitGroup
    wg.Add(1)
    go ms.StartGrpcServer(&wg)
  //  wg.Add(1)
   // go ms.StartHeartbeatChecker(&wg)
    wg.Add(1)
    go ms.StartDispatcher(&wg)
    wg.Add(1)
    go ms.StartServerForClient(&wg)
    wg.Wait()
}
