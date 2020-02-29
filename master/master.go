package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "os/signal"
    "sync"
    "syscall"
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

func exitByKill(ms *connect.Master, wg *sync.WaitGroup) {
    c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
    s := <- c
    log.Printf("Received Kill Signal, %v", s)
    //release
    ms.CloseClientConn()
    ms.ClearUpConnBeforeExit()
    wg.Done()
    os.Exit(0)
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
    wg.Add(1)
    go exitByKill(ms, &wg)
  //  wg.Add(1)
   // go ms.StartHeartbeatChecker(&wg)
    wg.Add(1)
    go ms.StartDispatcher(&wg)
    wg.Add(1)
    go ms.StartServerForClient(&wg)
    wg.Wait()
}
