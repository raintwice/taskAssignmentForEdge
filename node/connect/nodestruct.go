package connect

import (
	"google.golang.org/grpc"
	"sync"
	"taskAssignmentForEdge/taskmgt"
)

type Node struct {
	Maddr string
	Mport int
	Saddr string
	Sport int

	isOnline bool
	StartJoinTime int //duration to join after booted up

	MachineType int
	GroudIndex  int   //indicate one group with certain network environment
	BandWidth float64   //Unit: Mbps

	// the PstTime with the same group index should be the same
	PscTimeAvg int   //expectation of presence time in minutes
	PscTimeSigma int  //standard deviation in seconds
	Avl   float64  //availability

	conn *grpc.ClientConn  //call master
	//Tq *taskmgt.TaskQueue  //等待队列
	pool *taskmgt.Pool

	//pendingTasks []*taskmgt.TaskEntity  //等待传输时间
	//rwLock sync.RWMutex    //pendingTasks的锁

	PoolCap int //worker数量
	AvgExecTime int64  //平均任务执行时间
	CurTaskNum int  //队列长度即当前任务总数，等待队列长度= CurTaskNum - PoolCap if CurTaskNum < PoolCap else 0
	curTaskRwLock  sync.RWMutex//
	FinishTaskCnt int
}

func NewNode(addr string, port int, sport int) (* Node) {
	return &Node {
		Maddr:addr,
		Mport:port,
		Sport:sport,
		isOnline:false,
	}
}

func (no *Node) SetNodePara(bandWidth float64, machineType int, startTime int, poolCap int){
	no.BandWidth = bandWidth
	no.MachineType = machineType
	no.StartJoinTime = startTime
	no.PoolCap = poolCap
}

func (no *Node) SetNetworkPara(groupIndex int, pscTimeAvg int, pscTimeSig int, avl float64) {
	no.GroudIndex = groupIndex
	no.PscTimeAvg = pscTimeAvg
	no.PscTimeSigma = pscTimeSig
	no.Avl = avl
}

func  (no *Node) StartPool(cap int, wg *sync.WaitGroup) {
	no.pool = taskmgt.NewPool(cap)
	no.pool.Run()
	wg.Done()
}