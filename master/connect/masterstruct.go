package connect

import (
	"google.golang.org/grpc"
	"taskAssignmentForEdge/master/dispatch"
	"taskAssignmentForEdge/master/nodemgt"
	"taskAssignmentForEdge/taskmgt"
)

type Master struct {
	Tq *taskmgt.TaskQueue //全局等待队列
	Nq *nodemgt.NodeQueue //节点队列
	dispatcher dispatch.Dispatcher

	//
	ClientConn *grpc.ClientConn
}

func NewMaster() *Master {
	return &Master{
		Tq:   nil,
		Nq:   nil,
	}
}

func (ms *Master) Init() {
	ms.Tq = taskmgt.NewTaskQueue("global task queue")
	ms.Nq = nodemgt.NewNodeQueue()
	//TBD
	ms.dispatcher = dispatch.NewDefaultDispatcher()
}