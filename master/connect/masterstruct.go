package connect

import (
	"google.golang.org/grpc"
	"log"
	"taskAssignmentForEdge/master/dispatch"
	"taskAssignmentForEdge/master/nodemgt"
	"taskAssignmentForEdge/master/predictor"
	"taskAssignmentForEdge/taskmgt"
)

type Master struct {
	Tq *taskmgt.TaskQueue //全局等待队列
	Nq *nodemgt.NodeQueue //节点队列

	dispatcher dispatch.Dispatcher
	defaultDispachter dispatch.Dispatcher
	preDispatchCnt int
	PreDispatchNum int
	dispatchInterval int //in ms

	ClientConn *grpc.ClientConn
	runtimePredictMng *predictor.RunTimePredictManager
	connPredictMng *predictor.ConnPredictManager
}

func NewMaster() *Master {
	return &Master{
		Tq:   nil,
		Nq:   nil,
	}
}

func (ms *Master) Init(dispatchIndex int, dispatchItv int, preDispatchNum int) {
	ms.Tq = taskmgt.NewTaskQueue("global task queue")
	ms.Nq = nodemgt.NewNodeQueue()

	ms.defaultDispachter = dispatch.NewRRDispatcher()
	ms.dispatcher = dispatch.NewDispatcher(dispatchIndex)
	ms.runtimePredictMng = predictor.NewRunTimePredictManager()
	ms.connPredictMng = predictor.NewConnPredictManager()
	ms.dispatchInterval = dispatchItv
	ms.PreDispatchNum = preDispatchNum
	log.Printf("Pre-dispatch %d tasks", ms.PreDispatchNum)
}
