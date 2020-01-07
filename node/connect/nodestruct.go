package connect

import (
	"google.golang.org/grpc"
	"sync"
	"taskAssignmentForEdge/taskmgt"
)

const (
	MachineType_Simualted = iota
	MachineType_RaspPi_3B
	MachineType_RaspPi_4B
	MachineType_Surface_M3
)

type Node struct {
	Maddr string
	Mport int
	Saddr string
	Sport int

	isOnline bool
	StartJoinTime int //duration to join after booted up

	MachineType int
	PscTimeAvg int   //expectation of presence time in minutes
	PscTimeSigma int  //standard deviation in seconds
	Avl   float64  //availability


	conn *grpc.ClientConn  //与master的连接
	//Tq *taskmgt.TaskQueue  //等待队列
	pool *taskmgt.Pool

	//pendingTasks []*taskmgt.TaskEntity  //等待传输时间
	//rwLock sync.RWMutex    //pendingTasks的锁
	BandWidth float64   //Unit: MB/s
}

func NewNode(addr string, port int, sport int) (* Node) {
	return &Node {
		Maddr:addr,
		Mport:port,
		Sport:sport,
		isOnline:false,
	}
}

func (no *Node) SetParameters(bandWidth float64, machineType int, pscTimeAvg int, pscTimeSig int, avl float64, startTime int){
	no.BandWidth = bandWidth
	no.MachineType = machineType
	no.PscTimeAvg = pscTimeAvg
	no.PscTimeSigma = pscTimeSig
	no.Avl = avl
	no.StartJoinTime = startTime
}

func  (no *Node) StartPool(cap int, wg *sync.WaitGroup) {
	no.pool = taskmgt.NewPool(cap)
	no.pool.Run()
	wg.Done()
}