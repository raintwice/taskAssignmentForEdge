package connect

import (
	"google.golang.org/grpc"
	"log"
	"taskAssignmentForEdge/taskmgt"
)

type Node struct {
	Maddr string
	Mport int
	Saddr string
	Sport int

	isOnline bool
	StartJoinTime int //duration to join after booted up
	NodeMode int

	MachineType int
	GroudIndex  int   //indicate one group with certain network environment
	BandWidth float64   //Unit: Mbps
	Capacity float64

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
	FinishTaskCnt int

	/*
	CurTaskNum int  //队列长度即当前任务总数，等待队列长度= CurTaskNum - PoolCap if CurTaskNum < PoolCap else 0
	curTaskRwLock  sync.RWMutex//
	*/
	IsSetRecvServer bool
}

func NewNode(addr string, port int, sport int) (* Node) {
	return &Node {
		Maddr:addr,
		Mport:port,
		Sport:sport,
		isOnline:false,
		IsSetRecvServer:false,
	}
}

func (no *Node) SetNodePara(bandWidth float64, machineType int, startTime int, poolCap int, nodeMode int){
	no.BandWidth = bandWidth
	no.MachineType = machineType
	no.StartJoinTime = startTime
	no.PoolCap = poolCap
	no.NodeMode = nodeMode
}

func (no *Node) SetNetworkPara(groupIndex int, pscTimeAvg int, pscTimeSig int, avl float64) {
	no.GroudIndex = groupIndex
	no.PscTimeAvg = pscTimeAvg
	no.PscTimeSigma = pscTimeSig
	no.Avl = avl
}

/*
func  (no *Node) StartPool(cap int, wg *sync.WaitGroup) {
	no.pool = taskmgt.NewPool(cap)
	no.pool.Run()
	wg.Done()
}*/

func  (no *Node) StartPool(cap int) {
	no.pool = taskmgt.NewPool(cap)
	go no.pool.SafeRun()
}

func  (no *Node) StopPool() {
	if no.pool != nil {
		no.pool.SafeStop()
		no.pool = nil
	}
}

func (no *Node) SubmitTask(task *taskmgt.TaskEntity) {
	if no.pool != nil {
		no.pool.SafeSubmit(task)
	} else {
		log.Printf("Fatal:Task Pool is not booted up now in Node(%s:%d)", no.Saddr, no.Sport)
	}
}