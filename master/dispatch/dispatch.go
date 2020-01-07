package dispatch

import (
	"taskAssignmentForEdge/master/nodemgt"
	"taskAssignmentForEdge/taskmgt"
)

type Dispatcher interface {
	//EnqueueTask(tq *taskmgt.TaskQueue, task *taskmgt.TaskEntity )
	//EnqueueNode(nq *nodemgt.NodeQueue, node *nodemgt.NodeEntity )

	MakeDispatchDecision(tq *taskmgt.TaskQueue, nq *nodemgt.NodeQueue) (isNeedAssign bool)
	UpdateDispatcherStat(task *taskmgt.TaskEntity, node *nodemgt.NodeEntity)

	//Node(nq *nodemgt.NodeQueue, node *nodemgt.NodeEntity)
	//CheckNode(nq *nodemgt.NodeQueue)
}
