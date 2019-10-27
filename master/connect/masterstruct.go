package connect

import (
	"taskAssignmentForEdge/master/dispatch"
	"taskAssignmentForEdge/master/nodemgt"
	"taskAssignmentForEdge/taskmgt"
)

type Master struct {
	Tq *taskmgt.TaskQueue //等待队列
	Nq *nodemgt.NodeQueue //节点队列
	dispatcher dispatch.Dispatcher
}

func NewMaster() *Master {
	return &Master{
		Tq:   nil,
		Nq:   nil,
	}
}

func (ms *Master) Init() {
	ms.Tq = taskmgt.NewTaskQueue()
	ms.Nq = nodemgt.NewNodeQueue()
	ms.dispatcher = dispatch.NewDefaultDispatcher()
}