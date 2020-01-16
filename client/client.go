package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"taskAssignmentForEdge/client/connect"
)

var (
	h bool
	evalDir string
	arrvialRate float64
	pretrainNum int
	evalNum int
)

func init() {
	flag.BoolVar(&h, "h", false, "print help information")
	flag.StringVar(&evalDir, "dir", "Evalutionfiles", "set the directory of task files")
	flag.Float64Var(&arrvialRate, "rate", 150, "set the job arrival rate")
	flag.IntVar(&pretrainNum, "ptrainNum", 9000, "set the num of tasks used to pretrain tasks")
	flag.IntVar(&evalNum, "evalNum", 27000, "set the num of tasks used to evaluate ")
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: client [-h] [-dir DirofEvalutionfiles] [-rate arrvialRate] [-ptrainNum pretrainNum] [-evalNum evalNum]
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

	//log.Printf("rate:%f", arrvialRate)
	//log.Printf("pretrainNum:%d, evalNum:%d", pretrainNum, evalNum)
	client := connect.NewClient(evalDir, arrvialRate, pretrainNum, evalNum)
	client.InitConnection()
	var wg sync.WaitGroup
	wg.Add(1)
	go client.StartRecvResultServer(&wg)
	wg.Add(1)
	go client.StartTaskProducer(&wg)
	wg.Wait()
}